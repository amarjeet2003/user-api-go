package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amarjeet2003/user-api-go/repository"
	"github.com/amarjeet2003/user-api-go/routes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get database connection parameters from environment variables
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Construct the data source name
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
	log.Println("Connecting to database:", dataSourceName)

	// Connect to the database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Check if the connection is successful
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}

	log.Println("Connected to database successfully")

	// Check if the users table exists, create it if it doesn't
	err = ensureUsersTableExists(db)
	if err != nil {
		log.Fatal("Error ensuring users table exists:", err)
	}

	// Create a new UserRepository
	userRepo := repository.NewUserRepository(db)

	// Create a new ServeMux
	router := mux.NewRouter()

	// Setup user routes
	routes.SetupUserRoutes(router, userRepo)

	// Start the HTTP server
	log.Println("Server listening on :8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

// ensureUsersTableExists checks if the users table exists, creates it if it doesn't
func ensureUsersTableExists(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			first_name VARCHAR(50) NOT NULL,
			last_name VARCHAR(50) NOT NULL,
			username VARCHAR(50) UNIQUE NOT NULL,
			dob DATE NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	return nil
}
