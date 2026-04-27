package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/almanac12/e-library/book-service/internal/handler"
	"github.com/almanac12/e-library/book-service/internal/service"
)

func main(){
	r := gin.Default()

	svc := service.NewBookService()
	h := handler.NewBookHandler(svc)

	r.GET("/books", h.GetBooks)
	r.GET("/books/:id", h.GetBookByID)
	r.POST("/books", h.CreateBook)

	log.Println("server started on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}