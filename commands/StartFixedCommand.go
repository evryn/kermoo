package commands

import (
	"buggybox/modules/Router"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func GetStartFixedCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fixed",
		Short: "Use a fixed chance for mocking success or failure.",
		Example: `  Sleep for 8 seconds then start the service with the chance of %60 for successful responses:
    buggybox start fixed --chance 0.6 --sleep-for 8s
  Sleep for 1 minute then start with %20 success rate which is decided every 200ms
    buggybox start fixed --chance 0.2 --sleep-for 1m --interval 200ms
		`,
		Run: func(cmd *cobra.Command, args []string) {
			chance, _ := cmd.Flags().GetFloat32("success-chance")
			httpPort, _ := cmd.Flags().GetString("http-port")

			if chance < 0 || chance > 1 {
				fmt.Printf("Chance must be a float number between 0.0 and 1.0. Entered '%f' is not a valid value\n", chance)
				os.Exit(1)
			}

			Router.MustSetupRouter("0.0.0.0:" + httpPort)
		},
	}

	cmd.Flags().Float32P("success-chance", "c", 0.5, "The chance (0.0 to 1.0) of acting as a working app with successful response")

	return cmd
}
