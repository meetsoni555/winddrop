package cmd

import (
	"fmt"
	"os"
	"time"

	"winddrop/internal/file"
	"winddrop/internal/server"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("WindDrop CLI")
		fmt.Println("Usage:")
		fmt.Println("  winddrop send <files...> [--expire 5m] [--once]")
		return
	}

	command := os.Args[1]

	switch command {

	case "send":
		if len(os.Args) < 3 {
			fmt.Println("Please provide files or folders to send")
			return
		}

		var inputs []string
		var expiry time.Duration
		once := false

		for i := 2; i < len(os.Args); i++ {

			arg := os.Args[i]

			if arg == "--expire" && i+1 < len(os.Args) {
				dur, err := time.ParseDuration(os.Args[i+1])
				if err != nil {
					fmt.Println("Invalid duration")
					return
				}
				expiry = dur
				i++
				continue
			}

			if arg == "--once" {
				once = true
				continue
			}

			if _, err := os.Stat(arg); os.IsNotExist(err) {
				fmt.Println("Path does not exist:", arg)
				return
			}

			inputs = append(inputs, arg)
		}

		if len(inputs) == 0 {
			fmt.Println("No valid files provided")
			return
		}

		var fileToSend string
		isTempArchive := false

		if len(inputs) == 1 {
			info, _ := os.Stat(inputs[0])
			if info.IsDir() {
				fmt.Println("Creating archive...")
				archivePath, err := file.CreateArchive(inputs)
				if err != nil {
					fmt.Println("Archive failed:", err)
					return
				}
				fileToSend = archivePath
				isTempArchive = true
			} else {
				fileToSend = inputs[0]
			}
		} else {
			fmt.Println("Creating archive...")
			archivePath, err := file.CreateArchive(inputs)
			if err != nil {
				fmt.Println("Archive failed:", err)
				return
			}
			fileToSend = archivePath
			isTempArchive = true
		}

		server.StartServer(fileToSend, expiry, once, isTempArchive, len(inputs))

	default:
		fmt.Println("Unknown command:", command)
	}
}
