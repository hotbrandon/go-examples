package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	// Import the postgres driver
	_ "github.com/lib/pq"
)

// DB Configuration (Update these to match your local setup)
const (
	host   = "192.168.88.13"
	port   = 5433
	user   = "lisa"
	dbname = "postgres"
)

func main() {
	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		log.Fatal("DB_PASS environment variable not set")
	}
	password := dbPass

	// 1. Initialize Database Connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	fmt.Println("Successfully connected to database!")

	// 2. Start CLI Loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n--- Create User CLI ---\n")
		fmt.Print("Enter Username (or 'exit' to quit): ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		if username == "exit" {
			break
		}

		fmt.Print("Enter Phone: ")
		phone, _ := reader.ReadString('\n')
		phone = strings.TrimSpace(phone)

		// 3. Call the logic
		err := createUser(db, username, phone)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		} else {
			fmt.Println("✅ User created successfully!")
		}
	}
}

// --- Core Logic ---

// createUser manages the transaction lifecycle
func createUser(db *sql.DB, username, phone string) error {
	// 1. Start Transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %v", err)
	}
	// Safety net: Rollback ensures we don't leave hanging transactions on error.
	// If Commit() is called successfully later, this Rollback does nothing.
	defer tx.Rollback()

	// 2. Complex Check (Delegated to helper function)
	// We pass 'tx' so this check runs inside the current transaction scope
	exists, err := checkUserExists(tx, username)
	if err != nil {
		return fmt.Errorf("failed during existence check: %v", err)
	}

	if exists {
		return fmt.Errorf("username '%s' is already taken", username)
	}

	// 3. Insert New User
	insertQuery := "INSERT INTO users (username, phone) VALUES ($1, $2)"
	_, err = tx.Exec(insertQuery, username, phone)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	// 4. Commit
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("could not commit transaction: %v", err)
	}

	return nil
}

// checkUserExists encapsulates the "Complex Task"
// Note: It accepts *sql.Tx, not *sql.DB
func checkUserExists(tx *sql.Tx, username string) (bool, error) {
	// Imagine more complex logic here (e.g., checking archiving tables, external APIs, etc.)
	query := "SELECT username FROM users WHERE username = $1"

	var u string
	err := tx.QueryRow(query, username).Scan(&u)

	if err == sql.ErrNoRows {
		// No rows found = User does not exist
		return false, nil
	} else if err != nil {
		// Actual database error
		return false, err
	}

	// No error = Row found = User exists
	return true, nil
}
