package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

const (
	host     = "host"
	port     = "port"
	user     = "user"
	password = "password"
	dbname   = "postgres"
)

var db *sql.DB

func main() {
	var err error
	conn_str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", conn_str)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	create_table()

	router := mux.NewRouter()
	router.HandleFunc("/query", handle_query).Methods("GET")
	router.HandleFunc("/users", create_user).Methods("POST")
	router.HandleFunc("/users/{id}", get_user).Methods("GET")
	router.HandleFunc("/users/{id}", update_user).Methods("PUT")
	router.HandleFunc("/users/{id}", delete_user).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func create_table() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100),
			email VARCHAR(100) UNIQUE
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
}

func handle_query(w http.ResponseWriter, r *http.Request) {
	var request_body struct {
		Query string `json:"query"`
	}
	err := json.NewDecoder(r.Body).Decode(&request_body)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	rows, err := db.Query(request_body.Query)
	if err != nil {
		http.Error(w, "query execution failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, "failed to get column names", http.StatusInternalServerError)
		return
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		value_ptrs := make([]interface{}, len(columns))
		for i := range columns {
			value_ptrs[i] = &values[i]
		}
		if err := rows.Scan(value_ptrs...); err != nil {
			http.Error(w, "failed to scan row", http.StatusInternalServerError)
			return
		}
		row := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}
		result = append(result, row)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func create_user(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var id int
	err = db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&id)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func get_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var user struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "failed to get user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func update_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, id)
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func delete_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	rows_affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "failed to get rows affected", http.StatusInternalServerError)
		return
	}

	if rows_affected == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
