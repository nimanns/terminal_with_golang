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
	"time"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rand"
	"strings"
)

type item struct {
	id         int       `json:"id"`
	name       string    `json:"name"`
	created_at time.Time `json:"created_at"`
	updated_at time.Time `json:"updated_at"`
}

type user struct {
	id       int    `json:"id"`
	username string `json:"username"`
	password string `json:"password"`
}

var (
	items       = make(map[int]item)
	items_lock  sync.RWMutex
	next_item_id = 1

	users       = make(map[int]user)
	users_lock  sync.RWMutex
	next_user_id = 1

	sessions       = make(map[string]int)
	sessions_lock  sync.RWMutex
)

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	http.HandleFunc("/items", handle_items)
	http.HandleFunc("/items/", handle_item)
	http.HandleFunc("/upload", handle_upload)
	http.HandleFunc("/register", handle_register)
	http.HandleFunc("/login", handle_login)
	http.HandleFunc("/logout", handle_logout)
	http.HandleFunc("/search", handle_search)

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
	var new_item item
	err := json.NewDecoder(r.Body).Decode(&new_item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	items_lock.Lock()
	defer items_lock.Unlock()

	new_item.id = next_item_id
	next_item_id++
	new_item.created_at = time.Now()
	new_item.updated_at = time.Now()
	items[new_item.id] = new_item

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(new_item)
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
	updated_item.created_at = items[id].created_at
	updated_item.updated_at = time.Now()
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

func handle_register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var new_user user
	err := json.NewDecoder(r.Body).Decode(&new_user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users_lock.Lock()
	defer users_lock.Unlock()

	for _, u := range users {
		if u.username == new_user.username {
			http.Error(w, "username already exists", http.StatusConflict)
			return
		}
	}

	new_user.id = next_user_id
	next_user_id++
	new_user.password = hash_password(new_user.password)
	users[new_user.id] = new_user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(new_user)
}

func handle_login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var login_user user
	err := json.NewDecoder(r.Body).Decode(&login_user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	users_lock.RLock()
	defer users_lock.RUnlock()

	for _, u := range users {
		if u.username == login_user.username && u.password == hash_password(login_user.password) {
			session_token := generate_session_token()
			sessions_lock.Lock()
			sessions[session_token] = u.id
			sessions_lock.Unlock()

			http.SetCookie(w, &http.Cookie{
				Name:    "session_token",
				Value:   session_token,
				Expires: time.Now().Add(24 * time.Hour),
			})

			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "logged in successfully")
			return
		}
	}

	http.Error(w, "invalid credentials", http.StatusUnauthorized)
}

func handle_logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	session_token := c.Value

	sessions_lock.Lock()
	delete(sessions, session_token)
	sessions_lock.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "logged out successfully")
}

func handle_search(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "search query is required", http.StatusBadRequest)
		return
	}

	items_lock.RLock()
	defer items_lock.RUnlock()

	var results []item
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.name), strings.ToLower(query)) {
			results = append(results, item)
		}
	}

	json.NewEncoder(w).Encode(results)
}

func hash_password(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generate_session_token() string {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token)
}
