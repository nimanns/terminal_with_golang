package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

		"golang.org/x/crypto/bcrypt"
    _ "github.com/lib/pq"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "github.com/dgrijalva/jwt-go"
)

type config struct {
    db_host     string
    db_port     int
    db_user     string
    db_password string
    db_name     string
    jwt_secret  string
}

type user struct {
    id         int       `json:"id"`
    name       string    `json:"name"`
    email      string    `json:"email"`
    password   string    `json:"-"`
    role       string    `json:"role"`
    created_at time.Time `json:"created_at"`
    updated_at time.Time `json:"updated_at"`
}

var db *sql.DB
var conf config

func main() {
    var err error
    env_error := godotenv.Load()
    if env_error != nil {
        fmt.Println("error loading environment variables")
    }

    db_port, err_port := strconv.Atoi(os.Getenv("DBPORT"))
    if err_port != nil {
        fmt.Println("error converting port to int")
    }
    conf = config{
        db_host:     os.Getenv("DBHOST"),
        db_port:     db_port,
        db_user:     os.Getenv("DBUSER"),
        db_password: os.Getenv("DBPASSWORD"),
        db_name:     os.Getenv("DBNAME"),
        jwt_secret:  os.Getenv("JWT_SECRET"),
    }

    conn_str := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        conf.db_host, conf.db_port, conf.db_user, conf.db_password, conf.db_name)
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
    router.HandleFunc("/users", get_all_users).Methods("GET")
    router.HandleFunc("/users/{id}", get_user).Methods("GET")
    router.HandleFunc("/users/{id}", update_user).Methods("PUT")
    router.HandleFunc("/users/{id}", delete_user).Methods("DELETE")
    router.HandleFunc("/users/search", search_users).Methods("GET")
    router.HandleFunc("/login", login).Methods("POST")
    router.HandleFunc("/users/{id}/change_password", change_password).Methods("PUT")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }

    log.Printf("server starting on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, router))
}

func create_table() {
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100),
            email VARCHAR(100) UNIQUE,
            password VARCHAR(100),
            role VARCHAR(50) DEFAULT 'user',
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
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
    var new_user user
    err := json.NewDecoder(r.Body).Decode(&new_user)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    if new_user.role == "" {
        new_user.role = "user"
    }

    hashed_password, err := bcrypt.GenerateFromPassword([]byte(new_user.password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "failed to hash password", http.StatusInternalServerError)
        return
    }

    err = db.QueryRow("INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at", new_user.name, new_user.email, string(hashed_password), new_user.role).Scan(&new_user.id, &new_user.created_at, &new_user.updated_at)
    if err != nil {
        http.Error(w, "failed to create user", http.StatusInternalServerError)
        return
    }

    new_user.password = ""
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(new_user)
}

func get_all_users(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, email, role, created_at, updated_at FROM users")
    if err != nil {
        http.Error(w, "failed to get users", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []user
    for rows.Next() {
        var u user
        if err := rows.Scan(&u.id, &u.name, &u.email, &u.role, &u.created_at, &u.updated_at); err != nil {
            http.Error(w, "failed to scan user", http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func get_user(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "invalid user id", http.StatusBadRequest)
        return
    }

    var u user
    err = db.QueryRow("SELECT id, name, email, role, created_at, updated_at FROM users WHERE id = $1", id).Scan(&u.id, &u.name, &u.email, &u.role, &u.created_at, &u.updated_at)
    if err == sql.ErrNoRows {
        http.Error(w, "user not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "failed to get user", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(u)
}

func update_user(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "invalid user id", http.StatusBadRequest)
        return
    }

    var update_data user
    err = json.NewDecoder(r.Body).Decode(&update_data)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    _, err = db.Exec("UPDATE users SET name = $1, email = $2, role = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4", update_data.name, update_data.email, update_data.role, id)
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

func search_users(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("q")
    if query == "" {
        http.Error(w, "search query is required", http.StatusBadRequest)
        return
    }

    rows, err := db.Query("SELECT id, name, email, role, created_at, updated_at FROM users WHERE name ILIKE $1 OR email ILIKE $1", "%"+query+"%")
    if err != nil {
        http.Error(w, "failed to search users", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var users []user
    for rows.Next() {
        var u user
        if err := rows.Scan(&u.id, &u.name, &u.email, &u.role, &u.created_at, &u.updated_at); err != nil {
            http.Error(w, "failed to scan user", http.StatusInternalServerError)
            return
        }
        users = append(users, u)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func login(w http.ResponseWriter, r *http.Request) {
    var login_data struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    err := json.NewDecoder(r.Body).Decode(&login_data)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    var u user
    err = db.QueryRow("SELECT id, name, email, password, role FROM users WHERE email = $1", login_data.Email).Scan(&u.id, &u.name, &u.email, &u.password, &u.role)
    if err == sql.ErrNoRows {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    } else if err != nil {
        http.Error(w, "failed to get user", http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(u.password), []byte(login_data.Password))
    if err != nil {
        http.Error(w, "invalid credentials", http.StatusUnauthorized)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": u.id,
        "role":    u.role,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    token_string, err := token.SignedString([]byte(conf.jwt_secret))
    if err != nil {
        http.Error(w, "failed to generate token", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token_string})
}

func change_password(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "invalid user id", http.StatusBadRequest)
        return
    }

    var password_data struct {
        OldPassword string `json:"old_password"`
        NewPassword string `json:"new_password"`
    }

    err = json.NewDecoder(r.Body).Decode(&password_data)
    if err != nil {
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    var current_password string
    err = db.QueryRow("SELECT password FROM users WHERE id = $1", id).Scan(&current_password)
    if err == sql.ErrNoRows {
        http.Error(w, "user not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, "failed to get user", http.StatusInternalServerError)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(current_password), []byte(password_data.OldPassword))
    if err != nil {
        http.Error(w, "invalid old password", http.StatusUnauthorized)
        return
    }

    hashed_password, err := bcrypt.GenerateFromPassword([]byte(password_data.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "failed to hash password", http.StatusInternalServerError)
        return
    }

    _, err = db.Exec("UPDATE users SET password = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2", string(hashed_password), id)
    if err != nil {
        http.Error(w, "failed to update password", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
