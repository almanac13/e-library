package service

import (
	"borrow-service/internal/events"
	"borrow-service/internal/model"
	"borrow-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type BorrowService struct {
	repo      *repository.BorrowRepository
	publisher *events.Publisher
}

type CreateBorrowRequest struct {
	UserID string `json:"user_id" binding:"required"`
	BookID string `json:"book_id" binding:"required"`
}

func NewBorrowService(
	repo *repository.BorrowRepository,
	publisher *events.Publisher,
) *BorrowService {
	return &BorrowService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *BorrowService) CreateBorrow(req CreateBorrowRequest) (*model.Borrow, error) {
	now := time.Now()

	borrow := &model.Borrow{
		ID:         uuid.NewString(),
		UserID:     req.UserID,
		BookID:     req.BookID,
		BorrowDate: now,
		DueDate:    now.AddDate(0, 0, 14),
		Status:     model.BorrowStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := s.repo.Create(borrow)
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		s.publisher.Publish(
			"borrow.created",
			borrow,
		)
	}

	return borrow, nil
}

func (s *BorrowService) GetAllBorrows() ([]model.Borrow, error) {
	return s.repo.GetAll()
}

func (s *BorrowService) GetBorrowByID(id string) (*model.Borrow, error) {
	return s.repo.GetByID(id)
}

func (s *BorrowService) ReturnBorrow(id string) (*model.Borrow, error) {
	borrow, err := s.repo.ReturnBorrow(id, time.Now())
	if err != nil {
		return nil, err
	}

	if s.publisher != nil {
		s.publisher.Publish(
			"borrow.returned",
			borrow,
		)
	}

	return borrow, nil
}
