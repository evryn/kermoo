package commands

import (
	"buggybox/config"
	"fmt"

	"github.com/spf13/cobra"
)

func GetVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Get the BuggyBox version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(config.AppTitle)
			fmt.Println(config.AppDescription)
			fmt.Printf("Version: %s\nBuild: %s (%s)\n", config.BuildVersion, config.BuildRef, config.BuildDate)
			fmt.Println("Home: https://github.com/evryn/buggybox")
			fmt.Println("Source: https://github.com/evryn/buggybox")
			fmt.Println("Made with 💖 by: Amirreza Nasiri and contributors.")
		},
	}

	return cmd
}
