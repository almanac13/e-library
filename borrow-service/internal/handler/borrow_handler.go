package handler

import (
	"borrow-service/internal/service"
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BorrowHandler struct {
	service *service.BorrowService
}

func NewBorrowHandler(service *service.BorrowService) *BorrowHandler {
	return &BorrowHandler{service: service}
}

func (h *BorrowHandler) CreateBorrow(c *gin.Context) {
	var req service.CreateBorrowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	borrow, err := h.service.CreateBorrow(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, borrow)
}

func (h *BorrowHandler) GetAllBorrows(c *gin.Context) {
	borrows, err := h.service.GetAllBorrows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrows)
}

func (h *BorrowHandler) GetBorrowByID(c *gin.Context) {
	id := c.Param("id")

	borrow, err := h.service.GetBorrowByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "borrow not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrow)
}

func (h *BorrowHandler) ReturnBorrow(c *gin.Context) {
	id := c.Param("id")

	borrow, err := h.service.ReturnBorrow(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "borrow not found or already returned"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, borrow)
}
