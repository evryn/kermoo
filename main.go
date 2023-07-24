package main

import (
	"buggybox/modules/Datetime"
	"fmt"
)

// go run -ldflags "-X main.Version=1.0.0 -X 'main.Build=$(date --iso-8601)'" main.go

var Version string
var Build string

func main() {
	fmt.Println(`
	░░         
    ▒▒    ▓▓  ██
       ░░ ██    
  ▓▓    ▓▓        
    ▓▓▒▒   ░░   
	`)

	fmt.Printf("BuggyBox - %s (%s)\n", Version, Build)

	Datetime.MustSetInitialTime()
	fmt.Printf("initial time: %v\n", Datetime.InitialTime)
}
