package commands

import (
	"buggybox/modules/Time"
	"buggybox/modules/Utils"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func GetStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the BuggyBox foreground service",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			_, err := Utils.GetDurationFlag(cmd, "sleep-for")

			if err != nil {
				fmt.Printf("Invalid option: %s\n", err)
				os.Exit(1)
			}

			_, err = Utils.GetDurationFlag(cmd, "interval")

			if err != nil {
				fmt.Printf("Invalid option: %s\n", err)
				os.Exit(1)
			}

			port, _ := cmd.Flags().GetInt32("http-port")

			if port < 1 || port > 65535 {
				fmt.Printf("Invalid option: entered %d is not in valid port range (1-65535)\n", port)
				os.Exit(1)
			}

			port, _ = cmd.Flags().GetInt32("tcp-port")

			if port < 1 || port > 65535 {
				fmt.Printf("Invalid option: entered %d is not in valid port range (1-65535)\n", port)
				os.Exit(1)
			}
		},
	}

	cmd.PersistentFlags().StringP("sleep-for", "s", "5s", "Duration of delay before starting the service")
	cmd.PersistentFlags().StringP("interval", "i", "200ms", "Duration of the decision for each interval of success or failure")

	cmd.PersistentFlags().Int32("http-port", 80, "HTTP server port number to listen on")
	cmd.PersistentFlags().Int32("tcp-port", 5000, "TCP probe-able port number to listen on")
	cmd.PersistentFlags().String("state-dir", "/tmp/buggybox/", "Temporary directory to store state-related probe-able files")

	cmd.AddCommand(GetStartFixedCommand())

	return cmd
}

func InitStartPhase(cmd *cobra.Command) {
	sleepFor, _ := Utils.GetDurationFlag(cmd, "sleep-for")

	Time.MustSetInitialTime()

	fmt.Printf("Sleeping for %s...\n", sleepFor.String())
	time.Sleep(*sleepFor)
	fmt.Println("Waking up")
}

func InitWebServer(cmd *cobra.Command) {
	httpPort, _ := cmd.Flags().GetInt32("http-port")
	httpAddr := fmt.Sprintf("0.0.0.0:%d", httpPort)
	fmt.Printf("Starting HTTP server on %s\n", httpAddr)
	// web_server.MustSetupRouter(httpAddr)
}
