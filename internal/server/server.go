package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"winddrop/internal/utils"
)

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func StartServer(filePath string, expiry time.Duration, once bool, isTempArchive bool, totalItems int) {

	token := generateToken()

	fileName := filepath.Base(filePath)
	if isTempArchive {
		fileName = "winddrop_files.zip"
	}

	var expiryTime time.Time
	hasExpiry := expiry > 0

	if hasExpiry {
		expiryTime = time.Now().Add(expiry)
	}

	var downloaded int32 = 0

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 📥 DOWNLOAD HANDLER
	mux.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Query().Get("token") != token {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if hasExpiry && time.Now().After(expiryTime) {
			http.Error(w, "Link expired", http.StatusGone)
			return
		}

		if once && atomic.LoadInt32(&downloaded) == 1 {
			http.Error(w, "Already downloaded", http.StatusGone)
			return
		}

		mimeType := mime.TypeByExtension(filepath.Ext(fileName))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		w.Header().Set("Content-Type", mimeType)

		if once {
			atomic.StoreInt32(&downloaded, 1)
		}

		http.ServeFile(w, r, filePath)
	})

	ip := utils.GetLocalIP()

	// 🖥️ OUTPUT
	fmt.Println("\nWindDrop\n")

	if isTempArchive {
		fmt.Println("Mode      : Multi-file")
		fmt.Printf("Items     : %d\n", totalItems)
		fmt.Printf("Archive   : %s\n", fileName)
	} else {
		fmt.Printf("File      : %s\n", fileName)
	}

	if once {
		fmt.Println("Mode      : one-time")
	} else {
		fmt.Println("Mode      : normal")
	}

	if hasExpiry {
		fmt.Printf("Expires   : %v\n", expiry)
	} else {
		fmt.Println("Expires   : never")
	}

	fmt.Println()

	fmt.Printf("Link : http://%s:8080/download?token=%s\n", ip, token)

	fmt.Println("\nPress Ctrl+C to stop\n")

	// ⏳ CONTROL LOOP
	go func() {
		for {
			if hasExpiry {
				remaining := time.Until(expiryTime)
				if remaining <= 0 {
					fmt.Println("\nExpired. Shutting down...")
					server.Shutdown(context.Background())
					return
				}
				fmt.Printf("\r⏳ Time remaining: %-10v", remaining.Round(time.Second))
			}

			if once && atomic.LoadInt32(&downloaded) == 1 {
				fmt.Println("\nDownload complete. Shutting down...")
				time.Sleep(2 * time.Second)
				server.Shutdown(context.Background())
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()

	err := server.ListenAndServe()

	if isTempArchive {
		os.Remove(filePath)
	}

	if err != nil && err != http.ErrServerClosed {
		fmt.Println("Server error:", err)
	}
}
