package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/almanac12/e-library/book-service/internal/model"
)

type BookService struct{
	books map[string]model.Book
}

func NewBookService() *BookService{
	return &BookService{
		books: make(map[string]model.Book),
	}
}

func (s *BookService) GetAllBooks() []model.Book{
	result := make([]model.Book, 0, len(s.books))
	for _, b := range s.books {
		result = append(result, b)
	}
	return result
}

func (s *BookService) GetBookByID(id string) (model.Book, error){
	book, ok := s.books[id]
	if !ok{
		return model.Book{}, errors.New("not found")
	}
	return book, nil
}

func (s *BookService) CreateBook(b model.Book) model.Book {
	b.ID = uuid.New().String()
	b.Available = true
	s.books[b.ID] = b
	return b
}


