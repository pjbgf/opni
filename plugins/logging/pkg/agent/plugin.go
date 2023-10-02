package agent

import (
	"context"

	healthpkg "github.com/rancher/opni/pkg/health"
	"github.com/rancher/opni/pkg/logger"
	httpext "github.com/rancher/opni/pkg/plugins/apis/apiextensions/http"
	"github.com/rancher/opni/pkg/plugins/apis/apiextensions/stream"
	"github.com/rancher/opni/pkg/plugins/apis/capability"
	"github.com/rancher/opni/pkg/plugins/apis/health"
	"github.com/rancher/opni/pkg/plugins/driverutil"
	"github.com/rancher/opni/pkg/plugins/meta"
	"github.com/rancher/opni/plugins/logging/pkg/agent/drivers"
	"github.com/rancher/opni/plugins/logging/pkg/otel"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"go.uber.org/zap"
)

type Plugin struct {
	ctx           context.Context
	logger        *zap.SugaredLogger
	node          *LoggingNode
	otelForwarder *otel.Forwarder
}

func NewPlugin(ctx context.Context) *Plugin {
	lg := logger.NewPluginLogger().Named("logging")

	ct := healthpkg.NewDefaultConditionTracker(lg)

	p := &Plugin{
		ctx:    ctx,
		logger: lg,
		node:   NewLoggingNode(ct, lg),
		otelForwarder: otel.NewForwarder(
			otel.NewLogsForwarder(otel.WithLogger(lg.Named("otel-logs-forwarder"))),
			otel.NewTraceForwarder(otel.WithLogger(lg.Named("otel-trace-forwarder"))),
		),
	}

	for _, d := range drivers.NodeDrivers.List() {
		driverBuilder, _ := drivers.NodeDrivers.Get(d)
		driver, err := driverBuilder(ctx,
			driverutil.NewOption("logger", lg),
		)
		if err != nil {
			lg.With(
				"driver", d,
				zap.Error(err),
			).Error("failed to initialize logging node driver")
			continue
		}

		p.node.AddConfigListener(drivers.NewListenerFunc(ctx, driver.ConfigureNode))
	}
	return p
}

var (
	_ collogspb.LogsServiceServer   = (*otel.LogsForwarder)(nil)
	_ coltracepb.TraceServiceServer = (*otel.TraceForwarder)(nil)
)

func Scheme(ctx context.Context) meta.Scheme {
	scheme := meta.NewScheme(meta.WithMode(meta.ModeAgent))
	p := NewPlugin(ctx)
	scheme.Add(capability.CapabilityBackendPluginID, capability.NewAgentPlugin(p.node))
	scheme.Add(health.HealthPluginID, health.NewPlugin(p.node))
	scheme.Add(stream.StreamAPIExtensionPluginID, stream.NewAgentPlugin(p))
	scheme.Add(httpext.HTTPAPIExtensionPluginID, httpext.NewPlugin(p.otelForwarder))
	return scheme
}
