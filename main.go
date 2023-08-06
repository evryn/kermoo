package main

import (
	"buggybox/commands"
	"buggybox/config"
	"buggybox/modules/state"
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
	state.AppInitiatedAt = time.Now()
	commands.ExecuteRootCommand()

	//parseArgs()

	// fmt.Printf("Hostname %v\n", os.Getenv("HOSTNAME"))

	// Datetime.MustSetInitialTime()
	// fmt.Printf("Initialization time: %v\n", Datetime.InitialTime)

	// fmt.Printf("Sleeping for %s because of the initial delay...", InitDelay.String())
	// time.Sleep(*InitDelay)
	// fmt.Println("Woke up.")

}
