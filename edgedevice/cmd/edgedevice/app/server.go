package app

import (
	"fmt"
	"os"

	"github.com/kubeedge/beehive/pkg/core"
	"github.com/kubeedge/kubeedge/common/constants"
	"github.com/kubeedge/kubeedge/edgedevice/cmd/edgedevice/app/options"
	"github.com/kubeedge/kubeedge/edgedevice/pkg/common/dbm"
	"github.com/kubeedge/kubeedge/edgedevice/pkg/devicetwin"
	"github.com/kubeedge/kubeedge/edgedevice/pkg/edgehub"
	"github.com/kubeedge/kubeedge/edgedevice/pkg/eventbus"
	"github.com/kubeedge/kubeedge/pkg/features"
	"github.com/kubeedge/kubeedge/pkg/util"
	"github.com/kubeedge/kubeedge/pkg/util/flag"
	utilvalidation "github.com/kubeedge/kubeedge/pkg/util/validation"
	"github.com/kubeedge/kubeedge/pkg/version"
	"github.com/spf13/cobra"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/cli/globalflag"
	"k8s.io/component-base/term"
	"k8s.io/klog/v2"
)

// NewEdgeCoreCommand create edgedevice cmd
func NewEdgeDeviceCommand() *cobra.Command {
	opts := options.NewEdgeDeviceOptions()
	cmd := &cobra.Command{
		Use:  "edgedevice",
		Long: `Edgedevice is the component of device management on edge `,
		Run: func(cmd *cobra.Command, args []string) {
			flag.PrintMinConfigAndExitIfRequested(options.NewMinEdgeDeviceConfig())
			flag.PrintDefaultConfigAndExitIfRequested(options.NewDefaultEdgeDeviceConfig())
			flag.PrintFlags(cmd.Flags())

			if errs := opts.Validate(); len(errs) > 0 {
				klog.Exit(util.SpliceErrors(errs))
			}

			config, err := opts.Config()
			if err != nil {
				klog.Exit(err)
			}

			bootstrapFile := constants.BootstrapFile
			// get token from bootstrapFile if it exist
			if utilvalidation.FileIsExist(bootstrapFile) {
				token, err := os.ReadFile(bootstrapFile)
				if err != nil {
					klog.Exit(err)
				}
				config.Modules.EdgeHub.Token = string(token)
			}

			if err := features.DefaultMutableFeatureGate.SetFromMap(config.FeatureGates); err != nil {
				klog.Exit(err)
			}

			// To help debugging, immediately log version
			klog.Infof("Version: %+v", version.Get())

			registerModules(config)
			// start all modules
			core.Run()
		},
	}
	fs := cmd.Flags()
	namedFs := opts.Flags()
	flag.AddFlags(namedFs.FlagSet("global"))
	globalflag.AddGlobalFlags(namedFs.FlagSet("global"), cmd.Name())
	for _, f := range namedFs.FlagSets {
		fs.AddFlagSet(f)
	}

	usageFmt := "Usage:\n  %s\n"
	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStderr(), namedFs, cols)
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
		cliflag.PrintSections(cmd.OutOrStdout(), namedFs, cols)
	})

	return cmd
}

// registerModules register all the modules started in edgedevice
func registerModules(c *options.EdgeDeviceConfig) {
	hostnameOverride := util.GetHostname()

	devicetwin.Register(c.Modules.DeviceTwin, hostnameOverride)
	edgehub.Register(c.Modules.EdgeHub, hostnameOverride)
	eventbus.Register(c.Modules.EventBus, hostnameOverride)

	// database
	dbm.InitDBConfig(c.DataBase.DriverName, c.DataBase.AliasName, c.DataBase.DataSource)
}
