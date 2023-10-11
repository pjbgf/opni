package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"

	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	controlv1 "github.com/rancher/opni/pkg/apis/control/v1"
	corev1 "github.com/rancher/opni/pkg/apis/core/v1"
	"github.com/rancher/opni/pkg/config/v1beta1"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/plugins"
	"github.com/rancher/opni/pkg/storage"
	"github.com/rancher/opni/pkg/update"
	"github.com/rancher/opni/pkg/update/patch"
	"github.com/rancher/opni/pkg/urn"
	"github.com/spf13/afero"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FilesystemPluginSyncServer struct {
	controlv1.UnsafeUpdateSyncServer
	SyncServerOptions
	logger           *slog.Logger
	config           v1beta1.PluginsSpec
	loadMetadataOnce sync.Once
	manifest         *controlv1.UpdateManifest
	patchCache       patch.Cache
}

type SyncServerOptions struct {
	filters []plugins.Filter
	fsys    afero.Fs
}

type SyncServerOption func(*SyncServerOptions)

func (o *SyncServerOptions) apply(opts ...SyncServerOption) {
	for _, op := range opts {
		op(o)
	}
}

func WithPluginSyncFilters(filters plugins.Filter) SyncServerOption {
	return func(o *SyncServerOptions) {
		o.filters = append(o.filters, filters)
	}
}

func WithFs(fsys afero.Fs) SyncServerOption {
	return func(o *SyncServerOptions) {
		o.fsys = fsys
	}
}

func NewFilesystemPluginSyncServer(
	cfg v1beta1.PluginsSpec,
	lg *slog.Logger,
	opts ...SyncServerOption,
) (*FilesystemPluginSyncServer, error) {
	options := SyncServerOptions{
		fsys: afero.NewOsFs(),
	}
	options.apply(opts...)

	var patchEngine patch.BinaryPatcher
	switch cfg.Binary.Cache.PatchEngine {
	case v1beta1.PatchEngineBsdiff:
		patchEngine = patch.BsdiffPatcher{}
	case v1beta1.PatchEngineZstd:
		patchEngine = patch.ZstdPatcher{}
	default:
		return nil, fmt.Errorf("unknown patch engine: %s", cfg.Binary.Cache.PatchEngine)
	}

	var cache patch.Cache
	switch cfg.Binary.Cache.Backend {
	case v1beta1.CacheBackendFilesystem:
		var err error
		cache, err = patch.NewFilesystemCache(options.fsys, cfg.Binary.Cache.Filesystem, patchEngine, lg.WithGroup("cache"))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown cache backend: %s", cfg.Binary.Cache.Backend)
	}

	return &FilesystemPluginSyncServer{
		SyncServerOptions: options,
		config:            cfg,
		logger:            lg,
		patchCache:        cache,
	}, nil
}

func (f *FilesystemPluginSyncServer) Strategy() string {
	return patch.UpdateStrategy
}

func (f *FilesystemPluginSyncServer) RunGarbageCollection(ctx context.Context, store storage.ClusterStore) error {
	clusters, err := store.ListClusters(ctx, &corev1.LabelSelector{}, 0)
	if err != nil {
		return err
	}
	expected, err := f.CalculateExpectedManifest(ctx, urn.Plugin)
	if err != nil {
		return err
	}
	digestsToKeep := expected.DigestSet()
	for _, cluster := range clusters.Items {
		versions := cluster.GetMetadata().GetLastKnownConnectionDetails().GetPluginVersions()
		for _, v := range versions {
			digestsToKeep[v] = struct{}{}
		}
	}
	curDigests, err := f.patchCache.ListDigests()
	if err != nil {
		return err
	}
	var toClean []string
	for _, h := range curDigests {
		if _, ok := digestsToKeep[h]; !ok {
			toClean = append(toClean, h)
		}
	}
	f.logger.Info("running plugin cache gc")
	f.patchCache.Clean(toClean...)
	return nil
}

func (f *FilesystemPluginSyncServer) CalculateExpectedManifest(_ context.Context, updateType urn.UpdateType) (*controlv1.UpdateManifest, error) {
	if updateType != urn.Plugin {
		return nil, status.Error(codes.Unimplemented, fmt.Sprintf("unknown update type: %s", updateType))
	}
	f.loadMetadataOnce.Do(f.loadUpdateManifest)
	return f.manifest, nil
}

func (f *FilesystemPluginSyncServer) loadUpdateManifest() {
	if f.manifest != nil {
		panic("bug: tried to call loadUpdateManifest twice")
	}
	md, err := patch.GetFilesystemPlugins(plugins.DiscoveryConfig{
		Dir:        f.config.Dir,
		Fs:         f.fsys,
		Logger:     f.logger,
		Filters:    f.filters,
		QueryModes: len(f.filters) > 0,
	})
	if err != nil {
		panic(err)
	}
	if err := f.patchCache.Archive(md); err != nil {
		panic(fmt.Sprintf("failed to archive plugin manifest: %v", err))
	}
	f.manifest = md.ToManifest()
}

func (f *FilesystemPluginSyncServer) CalculateUpdate(
	ctx context.Context,
	theirManifest *controlv1.UpdateManifest,
) (*controlv1.PatchList, error) {
	// on startup
	if err := theirManifest.Validate(); err != nil {
		return nil, err
	}
	if items := theirManifest.GetItems(); len(items) > 0 {
		updateType, err := update.GetType(items)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		switch updateType {
		case urn.Plugin:
			return f.calculatePluginUpdates(ctx, theirManifest)
		default:
			return nil, status.Error(codes.Unimplemented, fmt.Sprintf("unknown update type: %s", updateType))
		}
	}
	return f.calculatePluginUpdates(ctx, theirManifest)
}

func (f *FilesystemPluginSyncServer) calculatePluginUpdates(
	ctx context.Context,
	theirManifest *controlv1.UpdateManifest,
) (*controlv1.PatchList, error) {
	ourManifest, err := f.CalculateExpectedManifest(ctx, urn.Plugin)
	if err != nil {
		return nil, err
	}
	archive := patch.LeftJoinOn(ourManifest, theirManifest)

	errg, _ := errgroup.WithContext(ctx)
	for _, entry := range archive.Items {
		entry := entry
		errg.Go(func() error {
			switch entry.Op {
			case controlv1.PatchOp_Create:
				data, err := f.patchCache.GetBinaryFile(patch.PluginsDir, entry.NewDigest)
				if err != nil {
					f.logger.Error("lost plugin in cache", logger.Err(err),
						"plugin", entry.Package,
						"filename", entry.Path)

					return status.Errorf(codes.Internal, "lost plugin in cache: %s", entry.Package)
				}
				entry.Data = data
			case controlv1.PatchOp_Update:
				// fetch existing patch or wait for a patch to be calculated
				lg := f.logger.With(
					"plugin", entry.Path,
					"oldDigest", entry.OldDigest,
					"newDigest", entry.NewDigest,
				)
				if data, err := f.patchCache.RequestPatch(entry.OldDigest, entry.NewDigest); err == nil {
					// send known patch
					entry.Data = data
				} else if errors.Is(err, os.ErrNotExist) {
					// no patch can ever be calculated in this case
					data, err := f.patchCache.GetBinaryFile(patch.PluginsDir, entry.NewDigest)
					if err != nil {
						lg.Error("lost plugin in cache", logger.Err(err))

						return status.Errorf(codes.Internal, "lost plugin in cache, cannot generate patch: %s", entry.Package)
					}
					entry.Data = data
					entry.Op = controlv1.PatchOp_Create
				} else {
					lg.Error(fmt.Sprintf(

						"error requesting patch for plugin %s %s->%s", entry.Package, entry.OldDigest, entry.NewDigest),
						logger.Err(err))

					return status.Errorf(codes.Internal, "internal error in plugin cache, cannot sync: %s", entry.Package)
				}
			}
			return nil
		})
	}
	if err := errg.Wait(); err != nil {
		return nil, err
	}

	return archive, nil
}

func (f *FilesystemPluginSyncServer) Collectors() []prometheus.Collector {
	return f.patchCache.MetricsCollectors()
}
