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
		fmt.Println("Winddrop CLI")
		fmt.Println("usage:")
		fmt.Println("  winddrop send <folder/files..> [--expire 5m] [--once] [--public]")
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
		var expiry time.Duration = 0
		once := false
		public := false

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

			if arg == "--public" {
				public = true
				continue
			}

			info, err := os.Stat(arg)
			if os.IsNotExist(err) {
				fmt.Println("Path does not exist:", arg)
				return
			}

			inputs = append(inputs, arg)
			_ = info
		}

		if len(inputs) == 0 {
			fmt.Println("No valid files/folders provided")
			return
		}

		var fileToSend string
		isTempArchive := false

		if len(inputs) == 1 {
			info, err := os.Stat(inputs[0])
			if err != nil {
				fmt.Println("Failed to read input:", err)
				return
			}

			if info.IsDir() {
				fmt.Println("Creating archive...")

				archivePath, err := file.CreateArchive(inputs)
				if err != nil {
					fmt.Println("Failed to create archive:", err)
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
				fmt.Println("Failed to create archive:", err)
				return
			}

			fileToSend = archivePath
			isTempArchive = true
		}

		fmt.Println("Starting WindDrop server...")

		server.StartServer(fileToSend, expiry, once, public, isTempArchive, len(inputs))

	default:
		fmt.Println("Unknown command:", command)
	}
}
