package main

import (
	"kermoo/commands"
	"kermoo/config"
	"time"
)

// go run main.go

var (
	Version string
	Hash    string
	Build   string

	InitDelay *time.Duration
)

func main() {

	config.InitializedAt = time.Now()
	commands.ExecuteRootCommand()

	//parseArgs()

	// fmt.Printf("Hostname %v\n", os.Getenv("HOSTNAME"))

	// Datetime.MustSetInitialTime()
	// fmt.Printf("Initialization time: %v\n", Datetime.InitialTime)

	// fmt.Printf("Sleeping for %s because of the initial delay...", InitDelay.String())
	// time.Sleep(*InitDelay)
	// fmt.Println("Woke up.")

}
