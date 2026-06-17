package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
)

var pingCount int64

func main() {
	filePath := "/usr/src/app/files/log.txt"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Write(data)
	})

	http.HandleFunc("/pings", func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt64(&pingCount, 1)
		fmt.Fprintf(w, "%d", count)
	})

	fmt.Println("Serving", filePath, "on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
