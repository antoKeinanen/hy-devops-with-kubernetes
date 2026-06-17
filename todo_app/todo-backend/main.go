package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Store struct {
	FilePath string
	Todos    []Todo
	mu       sync.Mutex
}

func NewStore(filePath string) *Store {
	s := &Store{FilePath: filePath}
	s.load()
	return s
}

func (s *Store) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Open(s.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.Todos = []Todo{}
			return
		}
		panic(err)
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	json.Unmarshal(bytes, &s.Todos)
}

func (s *Store) save() {
	bytes, _ := json.MarshalIndent(s.Todos, "", "  ")
	ioutil.WriteFile(s.FilePath, bytes, 0644)
}

func (s *Store) add(text string) Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := 1
	if len(s.Todos) > 0 {
		id = s.Todos[len(s.Todos)-1].ID + 1
	}

	todo := Todo{ID: id, Text: text}
	s.Todos = append(s.Todos, todo)
	s.save()

	return todo
}

func (s *Store) getAll() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.Todos
}

func main() {
	store := NewStore("todos.json")

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {

		case http.MethodGet:
			todos := store.getAll()
			json.NewEncoder(w).Encode(todos)

		case http.MethodPost:
			var text string

			if r.Header.Get("Content-Type") == "application/json" {
				var body map[string]string
				json.NewDecoder(r.Body).Decode(&body)
				text = body["text"]
			} else {
				r.ParseForm()
				text = r.FormValue("text")
			}

			if text == "" {
				http.Error(w, `{"error":"text is required"}`, http.StatusBadRequest)
				return
			}

			todo := store.add(text)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(todo)

		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	port := 8089
	fmt.Printf("Server running on http://localhost:%d\n", port)
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
