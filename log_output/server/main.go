package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Change this to the file you want to serve
	filePath := "/usr/src/app/files/log.txt"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filePath)
	})

	fmt.Println("Serving", filePath, "on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}

