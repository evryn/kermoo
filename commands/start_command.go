package commands

import (
	"buggybox/modules/logger"
	"buggybox/modules/user_config"
	"time"

	"github.com/spf13/cobra"
)

func GetStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the BuggyBox foreground service",
		Run: func(cmd *cobra.Command, args []string) {
			filename, _ := cmd.Flags().GetString("filename")
			verbosity, _ := cmd.Flags().GetString("verbosity")

			logger.MustInitLogger(verbosity)

			user_config.MustLoadProvidedConfig(filename)

			user_config.Prepared.Start()

			for {
				logger.Log.Info("app is alive")
				time.Sleep(1 * time.Minute)
			}
		},
	}

	cmd.Flags().StringP("filename", "f", "", "Path to your .yaml or .json configuration file")
	cmd.Flags().StringP("verbosity", "v", "", "Verbosity level of logging output. Valid values are: debug, info")

	return cmd
}
