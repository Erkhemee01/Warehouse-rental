package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID       int    // Add more fields as needed
	Phone    string
	LastName string
	SureName string
	Password string
}

type RegisterRequest struct {
	Phone    string `json:"phone"`
	LastName string `json:"lastName"`
	SureName string `json:"sureName"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:1234@tcp(localhost:3306)/Diploma")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping the database to verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")

	// Define a function to add CORS headers
	addCORSHeaders := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow requests from specific origin (replace '*' with the origin of your frontend application)
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			// If it's a preflight request, respond with a 200 status code
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Continue with the next handler
			next.ServeHTTP(w, r)
		})
	}

	// Create a new router
	router := http.NewServeMux()

	// Register handlers with the router
	router.HandleFunc("/register", RegisterHandler)
	router.HandleFunc("/login", LoginHandler)

	// Wrap the router with the CORS middleware
	handler := addCORSHeaders(router)

	// Start the server
	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
		return
	}

	// Validate request fields
	if req.Phone == "" || req.LastName == "" || req.SureName == "" || req.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Perform registration logic
	_, err := db.Exec("INSERT INTO users (phone, LastName, SureName, password) VALUES (?, ?, ?, ?)",
		req.Phone, req.LastName, req.SureName, req.Password)
	if err != nil {
		// Log the error for debugging
		log.Printf("Error inserting user data into database: %v", err)
		// Return a generic error message to the client
		http.Error(w, "Error inserting user data into database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Registration successful")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error parsing JSON data", http.StatusBadRequest)
		return
	}

	// Validate request fields
	if req.Phone == "" || req.Password == "" {
		http.Error(w, "Phone and password are required", http.StatusBadRequest)
		return
	}

	// Perform login logic
	var user User
	err := db.QueryRow("SELECT * FROM users WHERE phone = ? AND password = ?", req.Phone, req.Password).Scan(
		&user.ID, &user.Phone, &user.LastName, &user.SureName, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found or invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error querying user data from database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Login successful")
}