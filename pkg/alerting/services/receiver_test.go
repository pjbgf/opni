package services_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rancher/opni/internal/alertmanager"
	"github.com/rancher/opni/pkg/alerting/services"
	alertingv2 "github.com/rancher/opni/pkg/apis/alerting/v2"
	"github.com/rancher/opni/pkg/storage/inmemory"
	"github.com/rancher/opni/pkg/util"
	"github.com/samber/lo"
)

var _ = Describe("Receiver service", Label("unit"), Ordered, func() {
	var r services.ReceiverStorageService
	BeforeAll(func() {
		r = services.NewReceiverServer(
			inmemory.NewKeyValueStore[*alertingv2.OpniReceiver](
				func(in *alertingv2.OpniReceiver) *alertingv2.OpniReceiver { return util.ProtoClone(in) },
			),
		)
	})
	When("We use the receiver server", func() {
		It("should put a config", func() {
			ref, err := r.PutReceiver(context.TODO(), &alertingv2.OpniReceiver{
				Receiver: &alertmanager.Receiver{
					Name: lo.ToPtr("test"),
					SlackConfigs: []*alertmanager.SlackConfig{
						{
							ApiUrl: "https://slack.com/api",
						},
					},
				},
			})
			Expect(err).To(Succeed())
			Expect(ref).ToNot(BeNil())
		})

		It("should test a receiver", func() {
			// TODO : this requires ignoring config.Secret and amCfg.Secret in yaml Marshal / Unmarshal
			_, err := r.TestReceiver(context.TODO(), &alertingv2.OpniReceiver{
				Receiver: &alertmanager.Receiver{
					Name: lo.ToPtr("test"),
					SlackConfigs: []*alertmanager.SlackConfig{
						{
							ApiUrl: "https://slack.com/api",
						},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
