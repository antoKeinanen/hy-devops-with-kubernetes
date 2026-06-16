package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

var counter uint64

const outputFile = "/usr/src/app/files/log.txt"

func pingpongHandler(w http.ResponseWriter, r *http.Request) {
	count := atomic.AddUint64(&counter, 1)

	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	} else {
		defer file.Close()
		_, err := file.WriteString(fmt.Sprintf("ping / pong: %d\n", count))
		if err != nil {
			log.Printf("Error writing to file: %v", err)
		}
	}

	response := fmt.Sprintf("pong %d\n", count)
	w.Write([]byte(response))
}

func main() {
	http.HandleFunc("/pingpong", pingpongHandler)

	port := "8080"
	log.Printf("PingPong app running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
