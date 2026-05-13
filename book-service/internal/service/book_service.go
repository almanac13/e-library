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
func (s *BookService) Search(query string) ([]model.Book, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}

	return s.repo.Search(query)
}

func (s *BookService) ListByAuthor(author string) ([]model.Book, error) {
	if author == "" {
		return nil, errors.New("author is required")
	}

	return s.repo.FindByAuthor(author)
}

func (s *BookService) ListByCategory(category string) ([]model.Book, error) {
	if category == "" {
		return nil, errors.New("category is required")
	}

	return s.repo.FindByCategory(category)
}

func (s *BookService) CheckAvailability(id int) (bool, error) {
	book, err := s.GetByID(id)
	if err != nil {
		return false, err
	}

	return book.Available, nil
}

func (s *BookService) MarkAvailable(id int) (*model.Book, error) {
	book, err := s.repo.UpdateAvailability(id, true)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, errors.New("book not found")
	}

	return book, nil
}

func (s *BookService) MarkUnavailable(id int) (*model.Book, error) {
	book, err := s.repo.UpdateAvailability(id, false)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, errors.New("book not found")
	}

	return book, nil
}

func (s *BookService) GetStats() (*model.BookStats, error) {
	return s.repo.GetStats()
}
