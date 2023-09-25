package opensearchdata

import (
	"context"
	"fmt"

	"github.com/rancher/opni/pkg/util"
	"github.com/tidwall/gjson"
)

func (m *Manager) GetClusterStatus() ClusterStatus {
	if !m.IsInitialized() {
		return ClusterStatusNoClient
	}

	m.Lock()
	defer m.Unlock()

	resp, err := m.Client.Cluster.GetClusterHealth(context.TODO())
	if err != nil {
		m.logger.Error("failed to fetch opensearch cluster status", "err", err)
		return ClusterStatusError
	}
	defer resp.Body.Close()

	if resp.IsError() {
		m.logger.Error("failure response from cluster status", "resp", resp.String)
		return ClusterStatusError
	}

	respString := util.ReadString(resp.Body)
	status := gjson.Get(respString, "status").String()
	switch status {
	case "green":
		return ClusterStatusGreen
	case "yellow":
		return ClusterStatusYellow
	case "red":
		return ClusterStatusRed
	default:
		m.logger.Error(fmt.Sprintf("unknown status: %s", status))
		return ClusterStatusError
	}
}
