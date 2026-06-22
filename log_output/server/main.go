package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
)

var pingCount int64

func latestNonEmptyLine(content string) string {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return ""
	}

	lines := strings.Split(trimmed, "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

func main() {
	logFilePath := "/usr/src/app/files/log.txt"
	infoFilePath := "/usr/src/app/config/information.txt"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logData, err := os.ReadFile(logFilePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		infoData, err := os.ReadFile(infoFilePath)
		if err != nil {
			http.Error(w, "Information file not found", http.StatusInternalServerError)
			return
		}

		message := os.Getenv("MESSAGE")
		pingPongCount := atomic.LoadInt64(&pingCount)
		logLine := latestNonEmptyLine(string(logData))

		output := fmt.Sprintf(
			"file content: %s\nenv variable: MESSAGE=%s\n%s\nPing / Pongs: %d\n",
			strings.TrimSpace(string(infoData)),
			message,
			logLine,
			pingPongCount,
		)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Write([]byte(output))
	})

	http.HandleFunc("/pings", func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt64(&pingCount, 1)
		fmt.Fprintf(w, "%d", count)
	})

	fmt.Println("Serving", logFilePath, "on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
