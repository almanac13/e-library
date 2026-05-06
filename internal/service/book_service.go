package service

import (
	"errors"

	"github.com/almanac13/e-library/book-service/internal/model"
	"github.com/almanac13/e-library/book-service/internal/repository"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetAll() ([]model.Book, error) {
	return s.repo.GetAll()
}

func (s *BookService) GetByID(id int) (*model.Book, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, errors.New("book not found")
	}

	return book, nil
}

func (s *BookService) Create(req model.CreateBookRequest) (*model.Book, error) {
	if req.Title == "" || req.Author == "" || req.Category == "" {
		return nil, errors.New("title, author and category are required")
	}

	return s.repo.Create(req)
}

func (s *BookService) Update(id int, req model.UpdateBookRequest) (*model.Book, error) {
	if req.Title == "" || req.Author == "" || req.Category == "" {
		return nil, errors.New("title, author and category are required")
	}

	book, err := s.repo.Update(id, req)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, errors.New("book not found")
	}

	return book, nil
}

func (s *BookService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return errors.New("book not found")
	}

	return nil
}

func (s *BookService) FindByAuthor(author string) ([]model.Book, error) {
	if author == "" {
		return nil, errors.New("author query parameter is required")
	}

	return s.repo.FindByAuthor(author)
}
