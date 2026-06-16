package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

var counter uint64

func pingpongHandler(w http.ResponseWriter, r *http.Request) {
	count := atomic.AddUint64(&counter, 1)
	response := fmt.Sprintf("pong %d\n", count)
	w.Write([]byte(response))
}

func main() {
	http.HandleFunc("/pingpong", pingpongHandler)

	port := "8080"
	log.Printf("PingPong app running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
