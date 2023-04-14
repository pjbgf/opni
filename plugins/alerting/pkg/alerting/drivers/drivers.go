package drivers

import (
	"context"
	"sync"

	"github.com/rancher/opni/pkg/alerting/shared"
	corev1 "github.com/rancher/opni/pkg/apis/core/v1"
	"github.com/rancher/opni/plugins/alerting/pkg/apis/alertops"
	"github.com/samber/lo"
)

type ClusterDriverBuilder = func(ctx context.Context, args ...any) (ClusterDriver, error)

type ClusterDriver interface {
	alertops.AlertingAdminServer
	// Unique name of the driver
	Name() string
	// ShouldDisableNode is called during node sync for nodes which otherwise
	// have this capability enabled. If this function returns an error, the
	// node will be set to disabled instead, and the error will be logged.
	ShouldDisableNode(*corev1.Reference) error
	GetRuntimeOptions() shared.AlertingClusterOptions
}

var (
	lock                  = &sync.Mutex{}
	clusterDriverBuilders = make(map[string]ClusterDriverBuilder)
)

func RegisterClusterDriverBuilder(name string, fn ClusterDriverBuilder) {
	lock.Lock()
	defer lock.Unlock()

	clusterDriverBuilders[name] = fn
}

func UnregisterClusterDriverBuilder(name string) {
	lock.Lock()
	defer lock.Unlock()

	delete(clusterDriverBuilders, name)
}

func GetClusterDriverBuilder(name string) (ClusterDriverBuilder, bool) {
	lock.Lock()
	defer lock.Unlock()

	driver, ok := clusterDriverBuilders[name]
	return driver, ok
}

func ListClusterDrivers() []string {
	lock.Lock()
	defer lock.Unlock()

	return lo.Keys(clusterDriverBuilders)
}

type NoopClusterDriver struct {
	alertops.UnimplementedAlertingAdminServer
}

func (d *NoopClusterDriver) Name() string {
	return "noop"
}

func (d *NoopClusterDriver) ShouldDisableNode(*corev1.Reference) error {
	// the noop driver will never forcefully disable a node
	return nil
}

func (d *NoopClusterDriver) GetRuntimeOptions() shared.AlertingClusterOptions {
	return shared.AlertingClusterOptions{}
}

func init() {
	RegisterClusterDriverBuilder("noop", func(context.Context, ...any) (ClusterDriver, error) {
		return &NoopClusterDriver{}, nil
	})
}
