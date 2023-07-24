package commands

import (
	"buggybox/config"
	"fmt"

	"github.com/spf13/cobra"
)

func ExecuteRootCommand() {

	rootCommand := &cobra.Command{
		Use:   "buggybox",
		Short: config.AppDescription,
	}

	rootCommand.AddCommand(GetStartCommand())
	rootCommand.AddCommand(GetVersionCommand())
	rootCommand.Execute()
}

func printSignature() {
	fmt.Println(`
	     ░░         
       ▒▒    ▓▓  ██
          ░░ ██    
     ▓▓    ▓▓        
       ▓▓▒▒   ░░   
	`)

	fmt.Printf("%s %s\n", config.AppTitle, config.AppVersion)
}
