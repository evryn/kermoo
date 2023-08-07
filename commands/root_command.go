package commands

import (
	"buggybox/config"

	"github.com/spf13/cobra"
)

func ExecuteRootCommand() {

	rootCommand := &cobra.Command{
		Use:   "buggybox",
		Short: config.AppDescription,
	}

	rootCommand.AddCommand(GetStartCommand())
	rootCommand.AddCommand(GetVersionCommand())
	_ = rootCommand.Execute()
}

// func printSignature() {
// 	fmt.Println(`
// 	     ░░
//        ▒▒    ▓▓  ██
//           ░░ ██
//      ▓▓    ▓▓
//        ▓▓▒▒   ░░
// 	`)

// 	fmt.Printf("%s %s\n", config.AppTitle, config.BuildVersion)
// }
