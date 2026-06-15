package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func main() {
	id := uuid.New()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ts := time.Now().Format(time.RFC3339Nano)
		output := fmt.Sprintf("%s: %s\n", ts, id.String())

		fmt.Print(output)
		fmt.Fprint(w, output)
	})

	go func() {
		for {
			ts := time.Now().Format(time.RFC3339Nano)
			fmt.Printf("%s: %s\n", ts, id.String())
			time.Sleep(5 * time.Second)
		}
	}()

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
