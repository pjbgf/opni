package management_test

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rancher/opni-monitoring/pkg/core"
	"github.com/rancher/opni-monitoring/pkg/management"
	"github.com/rancher/opni-monitoring/pkg/storage"
	"github.com/rancher/opni-monitoring/pkg/test"
	"github.com/rancher/opni-monitoring/pkg/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Clusters", Ordered, Label(test.Unit, test.Slow), func() {
	var tv *testVars
	BeforeAll(setupManagementServer(&tv))

	It("should initially have no clusters", func() {
		clusters, err := tv.client.ListClusters(context.Background(), &management.ListClustersRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(clusters.Items).To(BeEmpty())
	})
	events := make(chan *management.WatchEvent, 1000)
	var streamCancel context.CancelFunc
	It("should handle watching create and delete events", func() {
		ctx, ca := context.WithCancel(context.Background())
		streamCancel = ca
		stream, err := tv.client.WatchClusters(ctx, &management.WatchClustersRequest{
			KnownClusters: &core.ReferenceList{
				Items: []*core.Reference{},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		go func() {
			defer close(events)
			for {
				event, err := stream.Recv()
				if err != nil {
					return
				}
				events <- event
			}
		}()
	})
	It("should create clusters", func() {
		for x := 0; x < 3; x++ {
			ids := map[string]struct{}{}
			for i := 0; i < 10; i++ {
				id := uuid.NewString()
				ids[id] = struct{}{}
				err := tv.storageBackend.CreateCluster(context.Background(), &core.Cluster{
					Id: id,
					Metadata: &core.ClusterMetadata{
						Labels: map[string]string{
							"i": fmt.Sprint(i + (x * 10)),
						},
					},
				})
				Expect(err).NotTo(HaveOccurred())
			}
			timeout := time.After(1100 * time.Millisecond)
			for i := 0; i < 10; i++ {
				select {
				case event := <-events:
					Expect(event.Type).To(Equal(management.WatchEventType_Added))
					Expect(ids).To(HaveKey(event.Cluster.Id))
					cluster, err := tv.client.GetCluster(context.Background(), &core.Reference{
						Id: event.Cluster.Id,
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(cluster.Metadata.Labels).To(HaveKey("i"))
					delete(ids, event.Cluster.Id)
				case <-timeout:
					Fail("timed out waiting for cluster create events")
				}
			}
			Expect(ids).To(BeEmpty())

			clusters, err := tv.client.ListClusters(context.Background(), &management.ListClustersRequest{})
			Expect(err).NotTo(HaveOccurred())
			Expect(clusters.Items).To(HaveLen(10 * (x + 1)))
		}
	})
	It("should edit cluster labels", func() {
		cluster, err := tv.client.ListClusters(context.Background(), &management.ListClustersRequest{
			MatchLabels: &core.LabelSelector{
				MatchExpressions: []*core.LabelSelectorRequirement{
					{
						Key:      "i",
						Operator: string(core.LabelSelectorOpIn),
						Values:   []string{"20"},
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(cluster.Items).To(HaveLen(1))
		ref := cluster.Items[0].Reference()
		updated, err := tv.client.EditCluster(context.Background(), &management.EditClusterRequest{
			Cluster: ref,
			Labels: map[string]string{
				"i": "999",
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(updated.Metadata.Labels).To(HaveKeyWithValue("i", "999"))

		updatedQueried, err := tv.client.GetCluster(context.Background(), ref)
		Expect(err).NotTo(HaveOccurred())
		Expect(updatedQueried.Metadata.Labels).To(HaveKeyWithValue("i", "999"))
	})
	It("should delete clusters", func() {
		clusters, err := tv.client.ListClusters(context.Background(), &management.ListClustersRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(clusters.Items).To(HaveLen(30))
		ids := map[string]struct{}{}
		for _, cluster := range clusters.Items {
			ids[cluster.Id] = struct{}{}
		}
		done := make(chan struct{})
		go func() {
			defer close(done)
			defer GinkgoRecover()

			for event := range events {
				Expect(event.Type).To(Equal(management.WatchEventType_Deleted))
				Expect(ids).To(HaveKey(event.Cluster.Id))
				delete(ids, event.Cluster.Id)

				_, err := tv.client.GetCluster(context.Background(), &core.Reference{
					Id: event.Cluster.Id,
				})
				Expect(status.Code(err)).To(Equal(codes.NotFound))
				if len(ids) == 0 {
					return
				}
			}
		}()
		for _, cluster := range clusters.Items {
			_, err := tv.client.DeleteCluster(context.Background(), &core.Reference{
				Id: cluster.Id,
			})
			Expect(err).NotTo(HaveOccurred())
			// watch events should be batched every second, wait 4 seconds in total
			// for all events to be received
			time.Sleep(100 * time.Millisecond)
		}
		Eventually(done).Should(BeClosed())

		clusters, err = tv.client.ListClusters(context.Background(), &management.ListClustersRequest{})
		Expect(err).NotTo(HaveOccurred())
		Expect(clusters.Items).To(BeEmpty())
		streamCancel()
	})
	When("attempting to edit a nonexistent cluster", func() {
		It("should error", func() {
			_, err := tv.client.EditCluster(context.Background(), &management.EditClusterRequest{
				Cluster: &core.Reference{
					Id: "nonexistent",
				},
				Labels: map[string]string{},
			})
			Expect(status.Code(err)).To(Equal(codes.NotFound))
		})
	})
	When("attempting to delete a nonexistent cluster", func() {
		It("should error", func() {
			_, err := tv.client.DeleteCluster(context.Background(), &core.Reference{
				Id: "nonexistent",
			})
			Expect(status.Code(err)).To(Equal(codes.NotFound))
		})
	})
	It("should handle validation errors", func() {
		_, err := tv.client.ListClusters(context.Background(), &management.ListClustersRequest{
			MatchLabels: &core.LabelSelector{
				MatchLabels: map[string]string{
					"\\": "bar",
				},
			},
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(validation.ErrInvalidLabelName.Error()))

		_, err = tv.client.GetCluster(context.Background(), &core.Reference{
			Id: "\\",
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(validation.ErrInvalidID.Error()))

		_, err = tv.client.EditCluster(context.Background(), &management.EditClusterRequest{
			Cluster: &core.Reference{
				Id: "\\",
			},
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(validation.ErrInvalidID.Error()))

		_, err = tv.client.DeleteCluster(context.Background(), &core.Reference{
			Id: "\\",
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(validation.ErrInvalidID.Error()))

		stream, err := tv.client.WatchClusters(context.Background(), &management.WatchClustersRequest{
			KnownClusters: &core.ReferenceList{
				Items: []*core.Reference{
					{
						Id: "\\",
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		_, err = stream.Recv()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(validation.ErrInvalidID.Error()))

		stream, err = tv.client.WatchClusters(context.Background(), &management.WatchClustersRequest{
			KnownClusters: &core.ReferenceList{
				Items: []*core.Reference{
					{
						Id: "nonexistent",
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		_, err = stream.Recv()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(storage.ErrNotFound.Error()))
	})
})
