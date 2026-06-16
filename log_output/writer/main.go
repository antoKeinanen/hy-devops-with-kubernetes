package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	id := uuid.New()

	filePath := "/usr/src/app/files/log.txt"

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fmt.Println("Writing to", filePath)

	for {
		ts := time.Now().Format(time.RFC3339Nano)
		output := fmt.Sprintf("%s: %s\n", ts, id.String())

		if _, err := file.WriteString(output); err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		fmt.Print(output)

		time.Sleep(5 * time.Second)
	}
}
