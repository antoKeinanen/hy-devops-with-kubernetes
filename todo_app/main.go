package main

import (
	"fmt"
	"net/http"
	"os"
)

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", hello)

	fmt.Printf("Server started in port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
