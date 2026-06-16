package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	imageURL     = "https://picsum.photos/1200"
	cacheDir     = "/usr/src/app/files"
	imagePath    = cacheDir + "/image.jpg"
	metadataPath = cacheDir + "/metadata.json"
	cacheTTL     = 10 * time.Minute
)

type Metadata struct {
	LastFetched time.Time `json:"last_fetched"`
}

var (
	mu           sync.Mutex
	serveOldOnce bool
)

func main() {
	os.MkdirAll(cacheDir, os.ModePerm)

	http.HandleFunc("/image", imageHandler)
	http.HandleFunc("/", indexHandler)

	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	content := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Todo App</title>
</head>

<body>
  <main class="container">
    <h1>Todo App</h1>

    <img
      class="image-card"
      src="/image"
      widht="250"
      height="250"
    />

    <form>
	<input placeholder="Enter a new todo (max 140 characters)" minlenght="1" maxlength="140" required />
	<input type="submit" />
    </form>

    <ul>
	<li>Learn kubernetes basics</li>
	<li>Deploy application to cluster</li>
	<li>Configure persistent groups</li>
    </ul>

    <p>DevOps with Kubernetes 2026</p>
  </main>
</body>
	</html>
	`
	w.Write([]byte(content))
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	meta := loadMetadata()

	if meta != nil && fileExists(imagePath) {
		expired := time.Since(meta.LastFetched) > cacheTTL

		if !expired {
			serveFile(w)
			return
		}

		if !serveOldOnce {
			serveOldOnce = true
			serveFile(w)
			return
		}
	}

	err := fetchAndStoreImage()
	if err != nil {
		http.Error(w, "Failed to fetch image", 500)
		return
	}

	serveOldOnce = false
	serveFile(w)
}

func fetchAndStoreImage() error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	meta := Metadata{
		LastFetched: time.Now(),
	}

	return saveMetadata(meta)
}

func serveFile(w http.ResponseWriter) {
	f, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, "No image available", 500)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	io.Copy(w, f)
}

func loadMetadata() *Metadata {
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil
	}

	var m Metadata
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}

	return &m
}

func saveMetadata(m Metadata) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return os.WriteFile(metadataPath, data, 0644)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
