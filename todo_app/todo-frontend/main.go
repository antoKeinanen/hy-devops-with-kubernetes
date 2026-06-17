package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"strings"
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

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
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
	todos, err := fetchTodos()
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}

	var todoItems strings.Builder
	for _, todo := range todos {
		todoItems.WriteString(fmt.Sprintf("<li>%s</li>", html.EscapeString(todo.Text)))
	}

	content := fmt.Sprintf(`
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
      width="250"
      height="250"
    />

    <form id="todoForm" action="/todos" method="POST">
      <input name="text" id="text" placeholder="Enter a new todo (max 140 characters)" minlength="1" maxlength="140" required />
      <input type="submit" />
    </form>

    <ul>
      %s
    </ul>

    <p>DevOps with Kubernetes 2026</p>
  </main>


<script>
  const form = document.getElementById("todoForm");

  form.addEventListener("submit", async function (event) {
    event.preventDefault();

    const text = document.getElementById("text").value;

    const response = await fetch("/todos", {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({ text })
    });

    if (!response.ok) {
      console.error("Failed to submit todo:", response.status);
      return;
    }

    window.location.reload();
  });

</script>

</body>
</html>
`, todoItems.String())

	w.Write([]byte(content))
}

func fetchTodos() ([]Todo, error) {
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8089"
	}

	resp, err := http.Get(backendURL + "/todos")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend returned status %d", resp.StatusCode)
	}

	var todos []Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, err
	}

	return todos, nil
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
