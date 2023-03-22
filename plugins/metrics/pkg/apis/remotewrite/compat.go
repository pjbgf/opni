package remotewrite

import (
	context "context"
	"sort"

	"github.com/cortexproject/cortex/pkg/cortexpb"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/prompb"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type PrometheusRemoteWriteClient interface {
	Push(ctx context.Context, in *prompb.WriteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type prometheusRemoteWriteClient struct {
	cc grpc.ClientConnInterface
}

func (c prometheusRemoteWriteClient) Push(ctx context.Context, in *prompb.WriteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/remotewrite.RemoteWrite/Push", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// This should not actually be used, but it is here to prevent mocks, tests,
// etc from breaking.
type slowPrometheusRemoteWriteClient struct {
	rwc RemoteWriteClient
}

func (s slowPrometheusRemoteWriteClient) Push(ctx context.Context, in *prompb.WriteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cwr := cortexpb.WriteRequest{}
	out, _ := in.Marshal()
	cwr.Unmarshal(out)
	_, err := s.rwc.Push(ctx, &cwr, opts...)
	if err != nil {
		return nil, err
	}
	cortexpb.ReuseSlice(cwr.Timeseries)
	return &emptypb.Empty{}, nil
}

func AsPrometheusRemoteWriteClient(rwc RemoteWriteClient) PrometheusRemoteWriteClient {
	switch rwc.(type) {
	case *remoteWriteClient:
		return prometheusRemoteWriteClient{cc: rwc.(*remoteWriteClient).cc}
	default:
		return slowPrometheusRemoteWriteClient{rwc: rwc}
	}
}

func labelProtosToLabels(labelPairs []prompb.Label) labels.Labels {
	result := make(labels.Labels, 0, len(labelPairs))
	for _, l := range labelPairs {
		result = append(result, labels.Label{
			Name:  l.Name,
			Value: l.Value,
		})
	}
	sort.Sort(result)
	return result
}
