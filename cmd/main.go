package main

import (
	"log"
	"net/http"
	"os"

	"github.com/almanac13/e-library/book-service/internal/db"
	"github.com/almanac13/e-library/book-service/internal/handler"
	"github.com/almanac13/e-library/book-service/internal/repository"
	"github.com/almanac13/e-library/book-service/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// Load variables from .env file.
	// Example: DB_HOST, DB_USER, DB_PASSWORD.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to PostgreSQL.
	database, err := db.ConnectPostgres()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	// Close database connection when app stops.
	defer database.Close()

	// Create repository.
	// Repository works with database.
	bookRepo := repository.NewBookRepository(database)
	bookService := service.NewBookService(bookRepo)
	// Create handler.
	// Handler works with HTTP requests.
	bookHandler := handler.NewBookHandler(bookService)

	// Route: /books
	// GET  /books  -> get all books
	// POST /books  -> create new book
	http.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookHandler.GetBooks(w, r)

		case http.MethodPost:
			bookHandler.CreateBook(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Route: /books/details?id=1
	// GET    -> get one book
	// PUT    -> update book
	// DELETE -> delete book
	http.HandleFunc("/books/details", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookHandler.GetBookByID(w, r)

		case http.MethodPut:
			bookHandler.UpdateBook(w, r)

		case http.MethodDelete:
			bookHandler.DeleteBook(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Route: /books/search?author=Martin
	http.HandleFunc("/books/search", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bookHandler.SearchBooks(w, r)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Read server port from .env.
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("book-service started on port", port)

	// Start HTTP server.
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
