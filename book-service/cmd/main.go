package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/almanac13/e-library/book-service/gen/bookpb"
	"github.com/almanac13/e-library/book-service/internal/cache"
	"github.com/almanac13/e-library/book-service/internal/db"
	"github.com/almanac13/e-library/book-service/internal/grpcserver"
	"github.com/almanac13/e-library/book-service/internal/handler"
	"github.com/almanac13/e-library/book-service/internal/repository"
	"github.com/almanac13/e-library/book-service/internal/service"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database, err := db.ConnectPostgres()
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer database.Close()

	if err := repository.RunMigration(database); err != nil {
		log.Fatal("failed to run migration:", err)
	}

	redisCache := cache.NewRedisCache()
	defer redisCache.Close()

	bookRepo := repository.NewBookRepository(database)
	bookService := service.NewBookService(bookRepo, redisCache)

	bookHandler := handler.NewBookHandler(bookService)
	bookGRPCServer := grpcserver.NewBookGRPCServer(bookService)

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

	http.HandleFunc("/books/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			bookHandler.SearchBooks(w, r)
			return
		}

		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal("failed to listen for gRPC:", err)
	}

	grpcServer := grpc.NewServer()
	bookpb.RegisterBookServiceServer(grpcServer, bookGRPCServer)

	go func() {
		log.Println("book-service gRPC started on port", grpcPort)

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("failed to serve gRPC:", err)
		}
	}()

	httpPort := os.Getenv("SERVER_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	log.Println("book-service HTTP started on port", httpPort)

	if err := http.ListenAndServe(":"+httpPort, nil); err != nil {
		log.Fatal(err)
	}
}
