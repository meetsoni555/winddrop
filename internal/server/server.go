package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"winddrop/internal/utils"
)

func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func StartServer(filePath string, expiryDuration time.Duration, once bool) {
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

	// ⏳ expiry countdown (only if expiry exists)
	if hasExpiry {
		go func() {
			for {
				remaining := time.Until(expiryTime)

				if remaining <= 0 {
					fmt.Println("\n⏳ Link expired. Shutting down WindDrop...")
					server.Shutdown(context.Background())
					return
				}

				fmt.Printf("\r⏳ Time remaining: %-10v", remaining.Round(time.Second))
				time.Sleep(1 * time.Second)
			}
		}()
	}

	mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {

		// token check
		queryToken := r.URL.Query().Get("token")
		if queryToken != token {
			http.Error(w, "❌ Unauthorized", http.StatusUnauthorized)
			return
		}

		// expiry check
		if hasExpiry && time.Now().After(expiryTime) {
			http.Error(w, "❌ Link expired", http.StatusGone)
			return
		}

		// one-time check
		if once && downloaded {
			http.Error(w, "❌ Already downloaded", http.StatusGone)
			return
		}

		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Type", "application/octet-stream")

		// mark as downloaded
		if once {
			downloaded = true
		}

		// shutdown if once mode
		if once {
			go func() {
				time.Sleep(2 * time.Second)
				fmt.Println("\n\n🛑 Transfer complete. WindDrop shutting down...")
				server.Shutdown(context.Background())
			}()
		}

		http.ServeFile(w, r, filePath)
	})

	ip := utils.GetLocalIP()

	// 🎨 Banner (clean version)
	fmt.Println(`
.................................................................
.................................................................
.................................................................
.##...##..######..##..##..#####...#####...#####....####...#####..
.##...##....##....###.##..##..##..##..##..##..##..##..##..##..##.
.##.#.##....##....##.###..##..##..##..##..#####...##..##..#####..
.#######....##....##..##..##..##..##..##..##..##..##..##..##.....
..##.##...######..##..##..#####...#####...##..##...####...##.....
.................................................................
.................................................................
.................................................................
`)

	fmt.Printf(" File      : %s\n", fileName)
	fmt.Printf(" Link      : http://%s:%s/download?token=%s\n", ip, port, token)
	fmt.Println(" Token     : enabled")

	if once {
		fmt.Println(" Mode      : one-time")
	} else {
		fmt.Println(" Mode      : normal")
	}

	// ✅ FIXED expiry logic
	if once && hasExpiry {
		fmt.Printf("⏳ Expires   : one-time or %v\n", expiryDuration)
	} else if once {
		fmt.Println("⏳ Expires   : one-time")
	} else if hasExpiry {
		fmt.Printf("⏳ Expires   : %v\n", expiryDuration)
	} else {
		fmt.Println("♾️ Expires   : never")
	}

	fmt.Println()
	fmt.Println("Press Ctrl+C to stop server")
	fmt.Println()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println("Server error:", err)
	}
}
