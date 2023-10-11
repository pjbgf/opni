package agent

import (
	"context"
	"fmt"

	"log/slog"

	healthpkg "github.com/rancher/opni/pkg/health"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/plugins/apis/apiextensions/stream"
	"github.com/rancher/opni/pkg/plugins/apis/capability"
	"github.com/rancher/opni/pkg/plugins/apis/health"
	"github.com/rancher/opni/pkg/plugins/meta"
	"github.com/rancher/opni/plugins/alerting/pkg/agent/drivers"
)

type Plugin struct {
	lg  *slog.Logger
	ctx context.Context

	ruleStreamer *RuleStreamer
	node         *AlertingNode
	driver       drivers.NodeDriver
}

func NewPlugin(ctx context.Context) *Plugin {
	lg := logger.NewPluginLogger().WithGroup("alerting")

	ct := healthpkg.NewDefaultConditionTracker(lg)
	p := &Plugin{
		ctx: ctx,
		lg:  lg,
	}

	p.node = NewAlertingNode(
		ctx,
		p.lg.With("component", "health-cfg-sync"),
		ct,
	)

	priority_order := []string{"k8s_driver", "test_driver"}
	for _, name := range priority_order {
		builder, ok := drivers.NodeDrivers.Get(name)
		if !ok {
			lg.Debug(fmt.Sprintf("could not find driver : %s", name))
			continue
		}
		driver, err := builder(ctx)
		if err != nil {
			lg.Warn("failed to initialize node driver", "driver", name, "err", err)
		}
		p.driver = driver
		p.node.AddConfigListener(driver)
		break
	}
	if p.driver == nil {
		panic("no driver set")
	}
	p.ruleStreamer = NewRuleStreamer(
		ctx,
		lg.With("component", "rule-streamer"),
		ct,
		p.driver,
	)
	p.node.AddConfigListener(p.ruleStreamer)
	return p
}

func Scheme(ctx context.Context) meta.Scheme {
	scheme := meta.NewScheme(meta.WithMode(meta.ModeAgent))
	p := NewPlugin(ctx)
	scheme.Add(capability.CapabilityBackendPluginID, capability.NewAgentPlugin(p.node))
	scheme.Add(health.HealthPluginID, health.NewPlugin(p.node))
	scheme.Add(stream.StreamAPIExtensionPluginID, stream.NewAgentPlugin(p))
	return scheme
}
