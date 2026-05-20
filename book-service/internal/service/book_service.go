package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/almanac13/e-library/book-service/internal/cache"
	"github.com/almanac13/e-library/book-service/internal/model"
	"github.com/almanac13/e-library/book-service/internal/repository"
	"github.com/redis/go-redis/v9"
)

type BookService struct {
	repo  *repository.BookRepository
	cache *cache.RedisCache
}

func NewBookService(repo *repository.BookRepository, redisCache *cache.RedisCache) *BookService {
	return &BookService{
		repo:  repo,
		cache: redisCache,
	}
}

func (s *BookService) GetAll() ([]model.Book, error) {
	ctx := context.Background()
	key := "books:all"

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, key)
		if err == nil {
			var books []model.Book
			if json.Unmarshal([]byte(cached), &books) == nil {
				return books, nil
			}
		}
	}

	books, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		data, _ := json.Marshal(books)
		_ = s.cache.Set(ctx, key, string(data), 5*time.Minute)
	}

	return books, nil
}

func (s *BookService) GetByID(id int) (*model.Book, error) {
	ctx := context.Background()
	key := fmt.Sprintf("book:%d", id)

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, key)
		if err == nil {
			var book model.Book
			if json.Unmarshal([]byte(cached), &book) == nil {
				return &book, nil
			}
		}

		if err != nil && !errors.Is(err, redis.Nil) {
			fmt.Println("redis get error:", err)
		}
	}

	book, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if book == nil {
		return nil, errors.New("book not found")
	}

	if s.cache != nil {
		data, _ := json.Marshal(book)
		_ = s.cache.Set(ctx, key, string(data), 5*time.Minute)
	}

	return book, nil
}

func (s *BookService) Create(req model.CreateBookRequest) (*model.Book, error) {
	if req.Title == "" || req.Author == "" || req.Category == "" {
		return nil, errors.New("title, author and category are required")
	}

	book, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	s.invalidateBooksCache(book.ID)

	return book, nil
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

	s.invalidateBooksCache(id)

	return book, nil
}

func (s *BookService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return errors.New("book not found")
	}

	s.invalidateBooksCache(id)

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

	s.invalidateBooksCache(id)

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

	s.invalidateBooksCache(id)

	return book, nil
}

func (s *BookService) GetStats() (*model.BookStats, error) {
	return s.repo.GetStats()
}

func (s *BookService) invalidateBooksCache(bookID int) {
	if s.cache == nil {
		return
	}

	ctx := context.Background()

	_ = s.cache.Delete(
		ctx,
		"books:all",
		fmt.Sprintf("book:%d", bookID),
	)
}
