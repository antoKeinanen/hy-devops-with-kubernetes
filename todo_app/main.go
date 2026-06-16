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

  <style>
    body {
      margin: 0;
      font-family: Arial, Helvetica, sans-serif;
      background: white;
      color: #333;
      text-align: center;
    }

    .container {
      padding-top: 85px;
    }

    h1 {
      font-size: 64px;
      margin: 0 0 45px;
      font-weight: 700;
    }

    .image-card {
      width: 400px;
      max-width: 90vw;
      height: 400px;
      object-fit: cover;
      border-radius: 14px;
      box-shadow: 0 4px 18px rgba(0, 0, 0, 0.25);
      display: block;
      margin: 0 auto;
    }

    .caption {
      margin-top: 45px;
      font-size: 32px;
      color: #666;
    }
  </style>
</head>

<body>
  <main class="container">
    <h1>Todo App</h1>

    <img
      class="image-card"
      src="/image"
    />

    <div class="caption">DevOps with Kubernetes 2026</div>
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
