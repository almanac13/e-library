package handler

import(
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/almanac12/e-library/book-service/internal/model"
	"github.com/almanac12/e-library/book-service/internal/service"
)

type BookHandler struct{
	svc *service.BookService
}

func NewBookHandler(svc *service.BookService) *BookHandler{
	return &BookHandler{svc: svc}
}

func (h *BookHandler) GetBooks(c *gin.Context){
	books := h.svc.GetAllBooks()
	c.JSON(http.StatusOK, books)
}

func (h *BookHandler) GetBookByID(c *gin.Context){
	id := c.Param("id")

	book, err := h.svc.GetBookByID(id)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandler) CreateBook(c *gin.Context){
	var req model.Book

	if err := c.ShouldBindJSON(&req); err !=nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	if req.Title == "" || req.Author == "" || req.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing fields"})
		return
	}

	book := h.svc.CreateBook(req)
	c.JSON(http.StatusCreated, book)
}
