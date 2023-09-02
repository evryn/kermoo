package commands

import (
	"kermoo/modules/logger"
	"kermoo/modules/user_config"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func GetStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [flags] [CONFIG]",
		Short: "Start the Kermoo foreground service",
		Long:  "Start the Kermoo foreground service. \n\nPass your config content or the config file path to the [CONFIG] argument below. You can also pass \"-\" to read from stdin. Leaving it empty will discover config from default path or `KERMOO_CONFIG` environmetn variable. \n\nRead documentation at: https://github.com/evryn/kermoo/wiki",
		Run: func(cmd *cobra.Command, args []string) {
			config, _ := cmd.Flags().GetString("filename")
			verbosity, _ := cmd.Flags().GetString("verbosity")

			if verbosity == "" {
				verbosity = os.Getenv("KERMOO_VERBOSITY")
			}

			logger.MustInitLogger(verbosity)

			if len(args) == 1 {
				config = args[0]
			}

			user_config.MustLoadPreparedConfig(config)

			user_config.Prepared.Start()

			for {
				logger.Log.Info("app is alive")
				time.Sleep(1 * time.Minute)
			}
		},
	}

	cmd.Flags().StringP("filename", "f", "", "(Deprecated) Content of config or path to file. Use [CONFIG] placeholder argument instead. It will be removed in future versions.")
	cmd.Flags().StringP("verbosity", "v", "", "Verbosity level of logging output, including: debug, info, warning, error, fatal. It overrides KERMOO_VERBOSITY environment variable.")

	return cmd
}
