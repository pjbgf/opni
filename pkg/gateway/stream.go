package gateway

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"log/slog"

	"github.com/kralicky/totem"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/metrics"
	"go.opentelemetry.io/otel/attribute"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	agentv1 "github.com/rancher/opni/pkg/agent"
	corev1 "github.com/rancher/opni/pkg/apis/core/v1"
	streamv1 "github.com/rancher/opni/pkg/apis/stream/v1"
	"github.com/rancher/opni/pkg/auth/cluster"
	"github.com/rancher/opni/pkg/plugins/meta"
	"github.com/rancher/opni/pkg/plugins/types"
	"github.com/rancher/opni/pkg/storage"
	"github.com/rancher/opni/pkg/util"
)

type streamPlugin struct {
	name string
	cc   *grpc.ClientConn
}

type internalRegistrar[T interface {
	registerInternalService(desc *grpc.ServiceDesc, impl any)
}] struct {
	source T
}

func (ir *internalRegistrar[T]) RegisterService(s *grpc.ServiceDesc, impl any) {
	ir.source.registerInternalService(s, impl)
}

type StreamServer struct {
	streamv1.UnimplementedStreamServer
	logger                   *slog.Logger
	handler                  ConnectionHandler
	clusterStore             storage.ClusterStore
	services                 []util.ServicePack[any]
	internalServices         []util.ServicePack[any]
	internalServiceRegistrar internalRegistrar[*StreamServer]
	streamPluginsMu          sync.Mutex
	streamPlugins            []streamPlugin
	metricsRegisterer        prometheus.Registerer

	providersMu  sync.Mutex
	providerById map[string]*metric.MeterProvider
}

func NewStreamServer(
	handler ConnectionHandler,
	clusterStore storage.ClusterStore,
	metricsRegisterer prometheus.Registerer,
	lg *slog.Logger,
) *StreamServer {
	srv := &StreamServer{
		logger:            lg.WithGroup("grpc"),
		handler:           handler,
		clusterStore:      clusterStore,
		metricsRegisterer: metricsRegisterer,
		providerById:      make(map[string]*metric.MeterProvider),
	}
	srv.internalServiceRegistrar.source = srv
	return srv
}

func (s *StreamServer) getProviderForId(agentId string) *metric.MeterProvider {
	s.providersMu.Lock()
	defer s.providersMu.Unlock()
	if prev, ok := s.providerById[agentId]; ok {
		return prev
	}
	exporter, err := otelprometheus.New(
		otelprometheus.WithRegisterer(prometheus.WrapRegistererWithPrefix("opni_gateway_", s.metricsRegisterer)),
		otelprometheus.WithoutScopeInfo(),
		otelprometheus.WithoutTargetInfo(),
	)
	if err != nil {
		s.logger.Error("failed to initialize stream metrics exporter", logger.Err(err))
		panic(err)
	}

	provider := metric.NewMeterProvider(metric.WithReader(exporter),
		metric.WithResource(resource.NewSchemaless(attribute.Key("agent-id").String(agentId))))
	s.providerById[agentId] = provider
	return provider
}

func (s *StreamServer) Connect(stream streamv1.Stream_ConnectServer) error {
	s.logger.Debug("handling new stream connection")
	ctx := stream.Context()

	id := cluster.StreamAuthorizedID(ctx)

	opts := []totem.ServerOption{
		totem.WithName("gateway"),
		totem.WithMetrics(s.getProviderForId(id),
			attribute.Key(metrics.LabelImpersonateAs).String(id),
		),
		totem.WithTracerOptions(
			resource.WithAttributes(
				semconv.ServiceNameKey.String("gateway"),
				attribute.String("agent", id),
			),
		),
	}

	ts, err := totem.NewServer(stream, opts...)
	if err != nil {
		return err
	}
	for _, service := range s.services {
		ts.RegisterService(service.Unpack())
	}

	c, err := s.clusterStore.GetCluster(ctx, &corev1.Reference{
		Id: id,
	})
	if err != nil {
		s.logger.Error("failed to get cluster", logger.Err(err),
			"id", id)

		return err
	}
	eventC, err := s.clusterStore.WatchCluster(ctx, c)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	ctx = storage.NewWatchContext(ctx, eventC)

	for _, r := range s.streamPlugins {
		streamClient := streamv1.NewStreamClient(r.cc)
		splicedStream, err := streamClient.Connect(ctx)
		name := fmt.Sprintf("gateway|%s", r.name)
		if err != nil {
			s.logger.Warn("failed to connect to remote stream, skipping", "clusterId", c.Id,
				logger.Err(err))

			continue
		}
		if err := ts.Splice(splicedStream,
			totem.WithName(name),
			totem.WithTracerOptions(resource.WithAttributes(
				semconv.ServiceNameKey.String(name),
				semconv.ServiceInstanceIDKey.String(id),
			)),
		); err != nil {
			s.logger.Warn("failed to splice remote stream, skipping", "clusterId", c.Id,
				logger.Err(err))

			continue
		}
	}

	cc, errC := ts.Serve()

	// check if an error was immediately returned
	select {
	case err := <-errC:
		return fmt.Errorf("stream connection failed: %w", err)
	default:
	}

	go s.handler.HandleAgentConnection(ctx, agentv1.NewClientSet(cc))

	select {
	case err = <-errC:
		if err != nil {
			s.logger.Warn("agent stream disconnected", logger.Err(err))

		}
		return status.Error(codes.Unavailable, err.Error())
	case <-ctx.Done():
		s.logger.Info("agent stream closing", logger.Err(ctx.Err()))

		err := ctx.Err()
		if errors.Is(err, storage.ErrObjectDeleted) {
			return status.Error(codes.Unauthenticated, err.Error())
		}
		return status.Error(codes.Unavailable, err.Error())
	}
}

func (s *StreamServer) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.logger.Debug("registering service", "service", desc.ServiceName)

	if len(desc.Streams) > 0 {
		s.logger.Error(

			"failed to register service: nested streams are currently not supported", "service", desc.ServiceName)
		panic("failed to register service: nested streams are currently not supported")

	}
	s.services = append(s.services, util.PackService(desc, impl))
}

func (s *StreamServer) registerInternalService(desc *grpc.ServiceDesc, impl any) {
	s.logger.Debug("registering internal service", "service", desc.ServiceName)

	if len(desc.Streams) > 0 {
		s.logger.Error(

			"failed to register internal service: nested streams are currently not supported", "service", desc.ServiceName)
		panic("failed to register internal service: nested streams are currently not supported")

	}
	s.internalServices = append(s.internalServices, util.PackService(desc, impl))
}

func (s *StreamServer) OnPluginLoad(ext types.StreamAPIExtensionPlugin, md meta.PluginMeta, cc *grpc.ClientConn) {
	lg := s.logger.With(
		"plugin", md.Filename(),
	)
	s.streamPluginsMu.Lock()
	defer s.streamPluginsMu.Unlock()
	lg.Debug("connecting to gateway plugin")
	s.streamPlugins = append(s.streamPlugins, streamPlugin{
		name: md.Filename(),
		cc:   cc,
	})

	internalStream, err := ext.ConnectInternal(context.Background())
	if err != nil {
		lg.Error("failed to connect to internal plugin stream", logger.Err(err))
		return
	}
	headerMd, err := internalStream.Header()
	if err != nil {
		lg.Error("failed to connect to internal plugin stream", logger.Err(err))
		return
	}
	var accepted bool
	if md := headerMd.Get("accept-internal-stream"); len(md) == 1 && md[0] == "true" {
		accepted = true
	}
	if !accepted {
		lg.Debug("plugin rejected internal stream connection")
		return
	}
	go func() {
		if err != nil {
			lg.Error("failed to connect to internal plugin stream", logger.Err(err))

			return
		}
		ts, err := totem.NewServer(
			internalStream,
			totem.WithName("gateway-internal-server"),
			totem.WithTracerOptions(
				resource.WithAttributes(
					semconv.ServiceNameKey.String("gateway-internal-server"),
					semconv.ServiceInstanceIDKey.String(md.Module),
				),
			),
		)
		if err != nil {
			lg.Error("failed to create internal plugin stream server", logger.Err(err))

			return
		}

		for _, service := range s.internalServices {
			ts.RegisterService(service.Unpack())
		}

		_, errC := ts.Serve()

		err = <-errC
		if err != nil {
			s.logger.Warn("internal plugin stream disconnected", logger.Err(err))

		}
	}()
}

func (s *StreamServer) InternalServiceRegistrar() grpc.ServiceRegistrar {
	return &s.internalServiceRegistrar
}
