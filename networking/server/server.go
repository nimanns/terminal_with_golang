package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type item struct {
	id   int    `json:"id"`
	name string `json:"name"`
}

var (
	items      = make(map[int]item)
	items_lock sync.RWMutex
	next_id    = 1
)

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/items", handle_items)
	http.HandleFunc("/items/", handle_item)

	http.HandleFunc("/upload", handle_upload)

	log.Println("server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handle_items(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		get_items(w, r)
	case http.MethodPost:
		create_item(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handle_item(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(filepath.Base(r.URL.Path))
	if err != nil {
		http.Error(w, "invalid item id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		get_item(w, r, id)
	case http.MethodPut:
		update_item(w, r, id)
	case http.MethodDelete:
		delete_item(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func get_items(w http.ResponseWriter, r *http.Request) {
	items_lock.RLock()
	defer items_lock.RUnlock()

	items_list := make([]item, 0, len(items))
	for _, item := range items {
		items_list = append(items_list, item)
	}

	json.NewEncoder(w).Encode(items_list)
}

func create_item(w http.ResponseWriter, r *http.Request) {
	var item item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items_lock.Lock()
	defer items_lock.Unlock()

	item.id = next_id
	next_id++
	items[item.id] = item

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func get_item(w http.ResponseWriter, r *http.Request, id int) {
	items_lock.RLock()
	defer items_lock.RUnlock()

	item, ok := items[id]
	if !ok {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func update_item(w http.ResponseWriter, r *http.Request, id int) {
	var updated_item item
	err := json.NewDecoder(r.Body).Decode(&updated_item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items_lock.Lock()
	defer items_lock.Unlock()

	if _, ok := items[id]; !ok {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	updated_item.id = id
	items[id] = updated_item

	json.NewEncoder(w).Encode(updated_item)
}

func delete_item(w http.ResponseWriter, r *http.Request, id int) {
	items_lock.Lock()
	defer items_lock.Unlock()

	if _, ok := items[id]; !ok {
		http.Error(w, "item not found", http.StatusNotFound)
		return
	}

	delete(items, id)
	w.WriteHeader(http.StatusNoContent)
}

func handle_upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dst, err := os.Create(fmt.Sprintf("./uploads/%s", header.Filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "file uploaded successfully: %s", header.Filename)
}
