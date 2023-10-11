//go:build !cli

package commands

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"syscall"

	"log/slog"

	"github.com/hashicorp/go-plugin"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"google.golang.org/grpc/codes"

	agentv2 "github.com/rancher/opni/pkg/agent/v2"
	"github.com/rancher/opni/pkg/bootstrap"
	"github.com/rancher/opni/pkg/config"
	"github.com/rancher/opni/pkg/config/v1beta1"
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/pkp"
	"github.com/rancher/opni/pkg/tokens"
	"github.com/rancher/opni/pkg/tracing"
	"github.com/rancher/opni/pkg/trust"
	"github.com/rancher/opni/pkg/util"

	_ "github.com/rancher/opni/pkg/ident/kubernetes"
	_ "github.com/rancher/opni/pkg/plugins/apis"
	_ "github.com/rancher/opni/pkg/storage/crds"
	_ "github.com/rancher/opni/pkg/storage/etcd"
	_ "github.com/rancher/opni/pkg/storage/jetstream"
	_ "github.com/rancher/opni/pkg/update/kubernetes/client"
	_ "github.com/rancher/opni/pkg/update/noop"
	_ "github.com/rancher/opni/pkg/update/patch/client"
)

func BuildAgentV2Cmd() *cobra.Command {
	var configFile, logLevel string
	var rebootstrap bool
	cmd := &cobra.Command{
		Use:   "agentv2",
		Short: "Run the v2 agent",
		Run: func(cmd *cobra.Command, args []string) {

			tracing.Configure("agentv2")
			agentlg := logger.New(logger.WithLogLevel(logger.ParseLevel(logLevel)))

			if configFile == "" {
				// find config file
				path, err := config.FindConfig()
				if err != nil {
					if errors.Is(err, config.ErrConfigNotFound) {
						wd, _ := os.Getwd()
						agentlg.Error(fmt.Sprintf(`could not find a config file in ["%s","/etc/opni"], and --config was not given`, wd))
						os.Exit(1)
					}
					agentlg.Error("an error occurred while searching for a config file", logger.Err(err))
					os.Exit(1)

				}
				agentlg.With(
					"path", path,
				).Info("using config file")
				configFile = path
			}

			objects, err := config.LoadObjectsFromFile(configFile)
			if err != nil {
				agentlg.Error("failed to load config", logger.Err(err))
				os.Exit(1)

			}
			var agentConfig *v1beta1.AgentConfig
			if ok := objects.Visit(func(config *v1beta1.AgentConfig) {
				agentConfig = config
			}); !ok {
				agentlg.Error("no agent config found in config file")
				os.Exit(1)
			}

			var bootstrapper bootstrap.Bootstrapper
			if agentConfig.Spec.ContainsBootstrapCredentials() {
				bootstrapper, err = configureBootstrapV2(agentConfig, agentlg)
				if err != nil {
					agentlg.Error("failed to configure bootstrap", logger.Err(err))
					os.Exit(1)

				}
			}

			p, err := agentv2.New(cmd.Context(), agentConfig,
				agentv2.WithBootstrapper(bootstrapper),
				agentv2.WithRebootstrap(rebootstrap),
			)
			if err != nil {
				agentlg.Error("error", logger.Err(err))
				return
			}

			err = p.ListenAndServe(cmd.Context())

			agentlg.Info("shutting down plugins")
			plugin.CleanupClients()
			agentlg.Info("all plugins shut down")

			if err != nil {
				const rebootstrapArg = "--re-bootstrap"
				var shouldRestart bool
				withoutArgs := []string{rebootstrapArg}
				var extraArgs []string

				if errors.Is(err, agentv2.ErrRebootstrap) {
					shouldRestart = true
					extraArgs = append(extraArgs, rebootstrapArg)
				} else if util.StatusCode(err) == codes.FailedPrecondition {
					shouldRestart = true
				}

				if shouldRestart {
					agentlg.Warn("preparing to restart agent", logger.Err(err))

					agentlg.Info(chalk.Yellow.Color("--- restarting agent ---"))
					args := append(lo.Without(os.Args, withoutArgs...), extraArgs...)
					panic(syscall.Exec(os.Args[0], args, os.Environ()))
				}
				if !errors.Is(err, context.Canceled) {
					agentlg.Error("error", logger.Err(err))
				}
			}
		},
	}
	cmd.Flags().StringVar(&configFile, "config", "", "Absolute path to a config file")
	cmd.Flags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warning, error)")
	cmd.Flags().BoolVar(&rebootstrap, "re-bootstrap", false, "attempt to re-bootstrap the agent even if it has already been bootstrapped")
	cmd.Flags().Lookup("re-bootstrap").Hidden = true
	return cmd
}

func configureBootstrapV2(conf *v1beta1.AgentConfig, agentlg *slog.Logger) (bootstrap.Bootstrapper, error) {
	var bootstrapper bootstrap.Bootstrapper
	var trustStrategy trust.Strategy
	if conf.Spec.Bootstrap == nil {
		return nil, errors.New("no bootstrap config provided")
	}
	if conf.Spec.Bootstrap.InClusterManagementAddress != nil {
		bootstrapper = &bootstrap.InClusterBootstrapperV2{
			GatewayEndpoint:    conf.Spec.GatewayAddress,
			ManagementEndpoint: *conf.Spec.Bootstrap.InClusterManagementAddress,
		}
	} else {
		agentlg.Info("loading bootstrap tokens from config file")
		tokenData := conf.Spec.Bootstrap.Token

		switch conf.Spec.TrustStrategy {
		case v1beta1.TrustStrategyPKP:
			var err error
			pins := conf.Spec.Bootstrap.Pins
			publicKeyPins := make([]*pkp.PublicKeyPin, len(pins))
			for i, pin := range pins {
				publicKeyPins[i], err = pkp.DecodePin(pin)
				if err != nil {
					agentlg.Error("failed to parse pin", "pin", string(pin), logger.Err(err))
					return nil, err
				}
			}
			conf := trust.StrategyConfig{
				PKP: &trust.PKPConfig{
					Pins: trust.NewPinSource(publicKeyPins),
				},
			}
			trustStrategy, err = conf.Build()
			if err != nil {
				agentlg.Error("error configuring PKP trust strategy", logger.Err(err))
				return nil, err
			}
		case v1beta1.TrustStrategyCACerts:
			paths := conf.Spec.Bootstrap.CACerts
			certs := []*x509.Certificate{}
			for _, path := range paths {
				data, err := os.ReadFile(path)
				if err != nil {
					agentlg.Error("failed to read CA cert", "path", path, logger.Err(err))
					return nil, err
				}
				cert, err := util.ParsePEMEncodedCert(data)
				if err != nil {
					agentlg.Error("failed to parse CA cert", "path", path, logger.Err(err))
					return nil, err
				}
				certs = append(certs, cert)
			}
			conf := trust.StrategyConfig{
				CACerts: &trust.CACertsConfig{
					CACerts: trust.NewCACertsSource(certs),
				},
			}
			var err error
			trustStrategy, err = conf.Build()
			if err != nil {
				agentlg.Error("error configuring CA Certs trust strategy", logger.Err(err))
				return nil, err
			}
		case v1beta1.TrustStrategyInsecure:
			agentlg.Warn(chalk.Bold.NewStyle().WithForeground(chalk.Yellow).Style(
				"*** Using insecure trust strategy. This is not recommended. ***",
			))
			conf := trust.StrategyConfig{
				Insecure: &trust.InsecureConfig{},
			}
			var err error
			trustStrategy, err = conf.Build()
			if err != nil {
				agentlg.Error("error configuring insecure trust strategy", logger.Err(err))
				return nil, err
			}
		}

		token, err := tokens.ParseHex(tokenData)
		if err != nil {
			agentlg.Error("failed to parse token", "token", fmt.Sprintf("[redacted (len: %d)]", len(tokenData)), logger.Err(err))
			return nil, err
		}
		bootstrapper = &bootstrap.ClientConfigV2{
			Token:         token,
			Endpoint:      conf.Spec.GatewayAddress,
			TrustStrategy: trustStrategy,
			FriendlyName:  conf.Spec.Bootstrap.FriendlyName,
		}
	}

	return bootstrapper, nil
}

func init() {
	AddCommandsToGroup(OpniComponents, BuildAgentV2Cmd())
}
