package main

import (
	"log"
	"net"
	"os"

	userpb "user-service/gen/userpb"
	"user-service/internal/db"
	grpcHandler "user-service/internal/grpc"
	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	pg, err := db.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	if err := repository.CreateUsersTable(pg); err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(pg)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	go startGRPCServer(userService)

	router := gin.Default()
	userHandler.RegisterRoutes(router)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8081"
	}

	log.Println("user-service HTTP started on port", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func startGRPCServer(userService *service.UserService) {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()

	userpb.RegisterUserServiceServer(
		server,
		grpcHandler.NewUserGRPCServer(userService),
	)

	reflection.Register(server)

	log.Println("user-service gRPC started on port", port)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
