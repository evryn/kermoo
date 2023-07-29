package commands

import (
	"buggybox/modules/logger"
	"buggybox/modules/user_config"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func GetStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the BuggyBox foreground service",
		Run: func(cmd *cobra.Command, args []string) {
			filename, _ := cmd.Flags().GetString("filename")
			verbosity, _ := cmd.Flags().GetString("verbosity")

			logger.MustInitLogger(verbosity)

			user_config.MustLoadUserConfig(filename)

			logger.Log.Info("configuration is loaded", zap.Any("configuration", user_config.UserConfig))

			user_config.UserConfig.Process.Run()
		},
	}

	cmd.Flags().StringP("filename", "f", "", "Path to your .yaml or .json configuration file")
	cmd.Flags().StringP("verbosity", "v", "", "Verbosity level of logging output. Valid values are: debug, info")

	return cmd
}
