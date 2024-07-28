package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
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

	http.HandleFunc("/query", handle_query)
	log.Fatal(http.ListenAndServe(":8081", nil))
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
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
