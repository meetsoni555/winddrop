package server

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"winddrop/internal/utils"
)

func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 🌍 Cloudflare tunnel
func startCloudflareTunnel() (string, *exec.Cmd) {
	cmd := exec.Command("cloudflared", "tunnel", "--url", "http://localhost:8080")

	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	err := cmd.Start()
	if err != nil {
		fmt.Println("❌ Failed to start cloudflared. Is it installed?")
		return "", nil
	}

	scanner := bufio.NewScanner(stdout)
	var publicURL string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "trycloudflare.com") {
			parts := strings.Fields(line)
			for _, p := range parts {
				if strings.HasPrefix(p, "https://") {
					publicURL = p
					break
				}
			}
			if publicURL != "" {
				break
			}
		}
	}

	return publicURL, cmd
}

func StartServer(filePath string, expiryDuration time.Duration, once bool, public bool) {
	port := "8080"

	token := generateToken()
	fileName := filepath.Base(filePath)

	var expiryTime time.Time
	hasExpiry := expiryDuration > 0

	if hasExpiry {
		expiryTime = time.Now().Add(expiryDuration)
	}

	downloaded := false

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// ⏳ expiry countdown
	if hasExpiry {
		go func() {
			for {
				remaining := time.Until(expiryTime)

				if remaining <= 0 {
					fmt.Println("\n⏳ Link expired. Shutting down...")
					server.Shutdown(context.Background())
					return
				}

				fmt.Printf("\r⏳ Time remaining: %-10v", remaining.Round(time.Second))
				time.Sleep(1 * time.Second)
			}
		}()
	}

	// 📥 download handler
	mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {

		// 🔐 token check
		queryToken := r.URL.Query().Get("token")
		if queryToken != token {
			http.Error(w, "❌ Unauthorized", http.StatusUnauthorized)
			return
		}

		// ⏳ expiry check
		if hasExpiry && time.Now().After(expiryTime) {
			http.Error(w, "❌ Link expired", http.StatusGone)
			return
		}

		// 🔥 one-time check
		if once && downloaded {
			http.Error(w, "❌ Already downloaded", http.StatusGone)
			return
		}

		// 📦 detect MIME properly
		mimeType := mime.TypeByExtension(filepath.Ext(fileName))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Type", mimeType)

		// mark downloaded
		if once {
			downloaded = true
		}

		// shutdown if once
		if once {
			go func() {
				time.Sleep(2 * time.Second)
				fmt.Println("\n🛑 Transfer complete. Shutting down...")
				server.Shutdown(context.Background())
			}()
		}

		http.ServeFile(w, r, filePath)
	})

	ip := utils.GetLocalIP()

	var publicURL string
	var tunnelCmd *exec.Cmd

	if public {
		fmt.Println("🌍 Starting Cloudflare tunnel...")
		publicURL, tunnelCmd = startCloudflareTunnel()

		if publicURL == "" {
			fmt.Println("❌ Failed to get public URL")
		}
	}

	// 🎨 output
	fmt.Println("\n🌬️ WindDrop\n")

	fmt.Printf("File      : %s\n", fileName)
	fmt.Println("Token     : enabled")

	if once {
		fmt.Println("Mode      : one-time")
	} else {
		fmt.Println("Mode      : normal")
	}

	if once && hasExpiry {
		fmt.Printf("Expires   : one-time or %v\n", expiryDuration)
	} else if once {
		fmt.Println("Expires   : one-time")
	} else if hasExpiry {
		fmt.Printf("Expires   : %v\n", expiryDuration)
	} else {
		fmt.Println("Expires   : never")
	}

	fmt.Println()

	// links
	fmt.Printf("Local Link  : http://%s:%s/download?token=%s\n", ip, port, token)

	if public && publicURL != "" {
		fmt.Printf("Public Link : %s/download?token=%s\n", publicURL, token)
	}

	fmt.Println("\nPress Ctrl+C to stop\n")

	err := server.ListenAndServe()

	// cleanup tunnel
	if tunnelCmd != nil && tunnelCmd.Process != nil {
		tunnelCmd.Process.Kill()
	}

	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Server error:", err)
	}
}
