package update

import (
	"context"
	"fmt"
	"time"

	controlv1 "github.com/rancher/opni/pkg/apis/control/v1"
	"github.com/rancher/opni/pkg/clients"
	"github.com/rancher/opni/pkg/util"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type SyncConfig struct {
	Client      controlv1.UpdateSyncClient
	StatsClient clients.ConnStatsQuerier
	Syncer      SyncHandler
	Logger      *zap.SugaredLogger
}

func (conf SyncConfig) DoSync(ctx context.Context) error {
	syncer := conf.Syncer
	initialManifest, err := getManifestWithTimeout(ctx, syncer)
	if err != nil {
		return err
	}
	updateType, err := GetType(initialManifest.GetItems())
	if err != nil {
		return fmt.Errorf("failed to get manifest update type: %w", err)
	}

	lg := conf.Logger.With(
		zap.String("type", string(updateType)),
	)
	lg.With(
		"entries", len(initialManifest.GetItems()),
	).Debug("sending manifest sync request")

	ticker := time.NewTicker(5 * time.Second)
	syncDone := make(chan struct{})
	statsDone := make(chan struct{})
	go func() {
		defer close(statsDone)
		var prevStats *clients.ConnStats
		startTime := time.Now()
		var printStats func(string)
		if conf.StatsClient == nil {
			var once bool
			printStats = func(string) {
				if once {
					return
				}
				once = true
			}
		} else {
			printStats = func(msg string) {
				stats, err := conf.StatsClient.QueryConnStats()
				if err != nil {
					lg.With(zap.Error(err)).Warn("failed to query connection stats")
					return
				}
				if prevStats == nil {
					prevStats = &stats
					return
				}
				_, rx := stats.CalcThroughput(*prevStats)
				rxStr, _ := util.Humanize(rx)
				recvdStr := stats.HumanizedBytesReceived()
				elapsedTime := time.Since(startTime)
				mins := elapsedTime / time.Minute
				elapsedTime -= mins * time.Minute
				secs := elapsedTime / time.Second
				elapsedTime -= secs * time.Second
				millis := elapsedTime / time.Millisecond
				lg.Debugf("%s%s | %sB/s | %02d:%02d.%03d", msg, recvdStr, rxStr, mins, secs, millis)
				prevStats = &stats
			}
		}
		printStats("")
		for {
			select {
			case <-ticker.C:
				printStats("retrieving patch data: ")
			case <-syncDone:
				printStats("patch data retrieved: ")
				return
			}
		}
	}()
	syncResp, err := conf.Client.SyncManifest(metadata.AppendToOutgoingContext(ctx,
		controlv1.UpdateStrategyKeyForType(updateType), syncer.Strategy(),
	), initialManifest)
	lg.Info("received sync response")
	syncDone <- struct{}{}
	ticker.Stop()
	<-statsDone
	if err != nil {
		return fmt.Errorf("failed to sync agent manifest: %w", err)
	}
	err = syncer.HandleSyncResults(ctx, syncResp)
	if err != nil {
		return fmt.Errorf("failed to handle agent sync results: %w", err)
	}
	lg.With(
		"entries", len(initialManifest.GetItems()),
	).Info("manifest sync complete")
	return nil
}

func (conf SyncConfig) Result(ctx context.Context) (*controlv1.UpdateManifest, error) {
	return getManifestWithTimeout(ctx, conf.Syncer)
}

func getManifestWithTimeout(ctx context.Context, syncer SyncHandler) (*controlv1.UpdateManifest, error) {
	ctx, ca := context.WithTimeout(ctx, 10*time.Second)
	m, err := syncer.GetCurrentManifest(ctx)
	ca()
	if err != nil {
		return nil, fmt.Errorf("failed to get current manifest: %w", err)
	}
	return m, nil
}
