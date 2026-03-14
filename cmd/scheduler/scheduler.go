package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marcosjom/sys-backups-automation/pkg/client"
	"github.com/marcosjom/sys-backups-automation/pkg/config"
)

func main() {

	configPath := ""
	run := false
	// Parse arguments
	argCount := len(os.Args)
	for i := 1; i < argCount; i++ {
		arg := os.Args[i]
		//log.Printf("Argument(%d / %d): '%s'\n", (i + 1), argCount, arg)
		switch arg {
		case "-c":
			if (i + i) >= argCount {
				fmt.Fprintf(os.Stderr, "Missing argument after: '%s'.\n", arg)
				return
			}
			if configPath != "" {
				fmt.Fprintf(os.Stderr, "Repeated argument: '%s'.\n", arg)
				return
			}
			configPath = os.Args[i+1]
			i++
		case "-r":
			if run {
				log.Printf("Repeated argument: '%s'.\n", arg)
			}
			run = true
		default:
			log.Printf("Ignoring argument: '%s'.\n", arg)
		}
	}
	//Help
	if argCount == 1 || configPath == "" {
		log.Printf("\n")
		log.Printf("This program executes periodic repetitive tasks under allowed windows of time.\n")
		log.Printf("\n")
		log.Printf("For example: \n")
		log.Printf(" - once every minute between its 5th-10th second.\n")
		log.Printf(" - once every day between 00:00-00:59.\n")
		log.Printf(" - twice a month the 1rst and 15th.\n")
		log.Printf(" - once a year during january.\n")
		log.Printf(" - etc...\n")
		log.Printf("\n")
		log.Printf("To verify a configuration file:\n")
		log.Printf("scheduler -c config.json\n")
		log.Printf("\n")
		log.Printf("To run:\n")
		log.Printf("scheduler -c config.json -r\n")
		return
	}
	// Load config
	if configPath != "" {
		config := config.Scheduler{}
		if err := config.Load(configPath); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: '%s': %s.\n", configPath, err.Error())
			return
		}
		log.Printf("Config loaded from: '%s'.\n", configPath)
		// Run client
		client := client.New()
		if err := client.Prepare(config.Client); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to prepare client with: '%s': %s.\n", configPath, err.Error())
			return
		}
		log.Printf("Client prepared.\n")
		// Run
		if run {
			log.Printf("Running...\n")
			tick := time.Tick(1 * time.Second)
			for {
				// Tick
				<-tick
				client.TickOneSecond()
			}
		}
		log.Printf("Done.\n")
	}
}
