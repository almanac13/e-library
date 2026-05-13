package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/almanac13/e-library/book-service/internal/model"
	"github.com/almanac13/e-library/book-service/internal/service"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetAll()
	if err != nil {
		http.Error(w, "failed to get books", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, books)
}

func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	book, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.service.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	var req model.UpdateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.service.Update(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromQuery(r)
	if err != nil {
		http.Error(w, "invalid book id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "book deleted successfully",
	})
}

func (h *BookHandler) SearchBooks(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	books, err := h.service.FindByAuthor(author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, books)
}

func getIDFromQuery(r *http.Request) (int, error) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		return 0, strconv.ErrSyntax
	}

	return strconv.Atoi(idParam)
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
