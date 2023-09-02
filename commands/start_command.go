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
		Use:   "start",
		Short: "Start the Kermoo foreground service",
		Run: func(cmd *cobra.Command, args []string) {
			config, _ := cmd.Flags().GetString("config")
			filename, _ := cmd.Flags().GetString("filename")
			verbosity, _ := cmd.Flags().GetString("verbosity")

			if verbosity == "" {
				verbosity = os.Getenv("KERMOO_VERBOSITY")
			}

			logger.MustInitLogger(verbosity)

			if filename != "" {
				logger.Log.Warn("`filename` flag is deprecated and will be removed in a future major release. use `config` flag instead.")
				config = filename
			}

			user_config.MustLoadPreparedConfig(config)

			user_config.Prepared.Start()

			for {
				logger.Log.Info("app is alive")
				time.Sleep(1 * time.Minute)
			}
		},
	}

	cmd.Flags().StringP("filename", "f", "", "Alias to `config` flag")
	cmd.Flags().StringP("config", "c", "", "Your YAML or JSON config content or path to a config file")
	cmd.Flags().StringP("verbosity", "v", "", "Verbosity level of logging output. Valid values are: debug, info")

	return cmd
}
