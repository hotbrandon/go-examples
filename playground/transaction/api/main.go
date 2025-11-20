package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// Imported to simulate delays if needed
	_ "github.com/lib/pq"
)

// --- 1. The Production Interface ---
// We renamed QueryRow -> QueryRowContext
// We added ExecContext so this interface can handle writes too if needed.
type Querier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// --- 2. The Shared Logic (Context Aware) ---

// checkUserExists now accepts 'ctx' as the first argument.
func checkUserExists(ctx context.Context, q Querier, username string) (bool, error) {
	query := "SELECT username FROM users WHERE username = $1"
	var u string

	// CRITICAL CHANGE: We use QueryRowContext and pass 'ctx'
	err := q.QueryRowContext(ctx, query, username).Scan(&u)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// --- 3. The HTTP Handlers ---

func main() {
	// Setup DB
	password := os.Getenv("DB_PASS")
	if password == "" {
		log.Fatal("DB_PASS environment variable not set")
	}

	connStr := fmt.Sprintf("host=192.168.88.13 port=5433 user=lisa password=%s dbname=postgres sslmode=disable", password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Using Go 1.22+ routing patterns
	http.HandleFunc("GET /check-user", handleCheckUser(db))
	http.HandleFunc("POST /create-user", handleCreateUser(db))

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCheckUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")

		// PASSING CONTEXT: r.Context() holds the request lifecycle
		exists, err := checkUserExists(r.Context(), db, username)

		if err != nil {
			// If the user cancelled the request, err will be "context canceled"
			log.Printf("Check error: %v", err)
			http.Error(w, "Internal Error", 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"exists": exists})
	}
}

func handleCreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		phone := "555-0000"
		ctx := r.Context() // Capture context once

		// 1. Start Transaction
		// Note: We use BeginTx now, not Begin!
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, "Tx Error", 500)
			return
		}
		defer tx.Rollback()

		// 2. Check (Passes Context + Transaction)
		exists, err := checkUserExists(ctx, tx, username)
		if err != nil {
			http.Error(w, "Check Error", 500)
			return
		}
		if exists {
			http.Error(w, "User exists", 409)
			return
		}

		// 3. Insert (Passes Context)
		// We use ExecContext
		_, err = tx.ExecContext(ctx, "INSERT INTO users (username, phone) VALUES ($1, $2)", username, phone)
		if err != nil {
			http.Error(w, "Insert Error", 500)
			return
		}

		// 4. Commit
		if err := tx.Commit(); err != nil {
			http.Error(w, "Commit Error", 500)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "User Created")
	}
}
