package main

import (
	"buggybox/modules/Datetime"
	"flag"
	"fmt"
	"os"
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
	printSignature()
	parseArgs()

	fmt.Printf("Hostname %v\n", os.Getenv("HOSTNAME"))

	Datetime.MustSetInitialTime()
	fmt.Printf("Initialization time: %v\n", Datetime.InitialTime)

	fmt.Printf("Sleeping for %s because of the initial delay...", InitDelay.String())
	time.Sleep(*InitDelay)
	fmt.Println("Woke up.")

}

func printSignature() {
	fmt.Println(`
	     ░░         
       ▒▒    ▓▓  ██
          ░░ ██    
     ▓▓    ▓▓        
       ▓▓▒▒   ░░   
	`)

	fmt.Printf("BuggyBox %s - %s (%s)\n", Version, Hash, Build)

}

func parseArgs() {
	InitDelay = flag.Duration("init-delay", time.Duration(5*time.Second), "Initial container delay")
	flag.Parse()
}
