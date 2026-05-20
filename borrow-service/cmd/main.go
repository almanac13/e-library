package main

import (
	"log"
	"net"
	"os"

	"borrow-service/gen/borrowpb"
	"borrow-service/internal/db"
	"borrow-service/internal/events"
	borrowgrpc "borrow-service/internal/grpc"
	"borrow-service/internal/handler"
	"borrow-service/internal/repository"
	"borrow-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	database := db.Connect()

	err := repository.RunMigration(database)
	if err != nil {
		log.Fatal(err)
	}

	borrowRepo := repository.NewBorrowRepository(database)

	publisher := events.NewPublisher()
	defer publisher.Close()

	borrowService := service.NewBorrowService(
		borrowRepo,
		publisher,
	)

	borrowHandler := handler.NewBorrowHandler(borrowService)

	go startGRPCServer(borrowService)

	router := gin.Default()

	router.POST("/borrows", borrowHandler.CreateBorrow)
	router.GET("/borrows", borrowHandler.GetAllBorrows)
	router.GET("/borrows/:id", borrowHandler.GetBorrowByID)
	router.PUT("/borrows/:id/return", borrowHandler.ReturnBorrow)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8089"
	}

	log.Println("borrow-service HTTP started on port " + port)

	err = router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}

func startGRPCServer(borrowService *service.BorrowService) {
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50053"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	borrowpb.RegisterBorrowServiceServer(
		grpcServer,
		borrowgrpc.NewBorrowGRPCServer(borrowService),
	)

	log.Println("borrow-service gRPC started on port " + port)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal(err)
	}
}
