package cmd

import (
	"fmt"
	"os"
	"time"

	"winddrop/internal/server"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("WindDrop CLI")
		fmt.Println("Usage:")
		fmt.Println("  winddrop send <file> [--expire 5m] [--once]")
		return
	}

	command := os.Args[1]

	switch command {

	case "send":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a file to send")
			return
		}

		file := os.Args[2]

		// check file exists
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Println("❌ File does not exist:", file)
			return
		}

		var expiry time.Duration = 0
		once := false

		// parse flags
		for i := 3; i < len(os.Args); i++ {

			if os.Args[i] == "--expire" && i+1 < len(os.Args) {
				dur, err := time.ParseDuration(os.Args[i+1])
				if err != nil {
					fmt.Println("❌ Invalid duration (use 2m, 10m)")
					return
				}
				expiry = dur
				i++
			}

			if os.Args[i] == "--once" {
				once = true
			}
		}

		fmt.Println("Starting WindDrop server for:", file)

		server.StartServer(file, expiry, once)

	default:
		fmt.Println("Unknown command:", command)
	}
}
