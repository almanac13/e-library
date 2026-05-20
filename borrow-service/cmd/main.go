package main

import (
	"log"
	"os"

	"borrow-service/internal/db"
	"borrow-service/internal/handler"
	"borrow-service/internal/repository"
	"borrow-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	database := db.Connect()

	err = repository.RunMigration(database)
	if err != nil {
		log.Fatal(err)
	}

	borrowRepo := repository.NewBorrowRepository(database)
	borrowService := service.NewBorrowService(borrowRepo)
	borrowHandler := handler.NewBorrowHandler(borrowService)

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
