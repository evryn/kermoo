package commands

import (
	"kermoo/config"

	"github.com/spf13/cobra"
)

func ExecuteRootCommand() {

	rootCommand := &cobra.Command{
		Use:   "kermoo",
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
