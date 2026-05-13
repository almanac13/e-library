package main

import (
	"log"
	"os"

	"user-service/internal/db"
	"user-service/internal/handler"
	"user-service/internal/repository"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	pg, err := db.NewPostgresConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	err = repository.CreateUsersTable(pg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(pg)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

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
