package main

import (
	"fmt"
	"os"
)

func main() {

	configPath := ""
	run := false
	// Parse arguments
	iArg := 1
	argCount := len(os.Args)
	for i := 1; i < argCount; i++ {
		arg := os.Args[i]
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
				fmt.Printf("Repeated argument: '%s'.\n", arg)
			}
			run = true
		default:
			fmt.Printf("Ignoring argument: '%s'.\n", arg)
		}
	}
	//Help
	if iArg == 1 || configPath == "" {
		fmt.Printf("\n")
		fmt.Printf("This program execute periodic tasks under allowed window-time.\n")
		fmt.Printf("Intended for database and/or files backups but universal.\n")
		fmt.Printf("\n")
		fmt.Printf("Usage:\n")
		fmt.Printf("scheduler -c config.json -r\n")
		fmt.Printf("\n")
		fmt.Printf("Examples:\n")
		fmt.Printf("scheduler -c config.json\n")
		fmt.Printf("scheduler -c config.json -r\n")
	}

}
