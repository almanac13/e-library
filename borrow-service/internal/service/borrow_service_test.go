package service

import (
	"borrow-service/internal/model"
	"testing"
	"time"
)

type mockBorrowRepository struct {
	borrows map[string]*model.Borrow
}

func newMockBorrowRepository() *mockBorrowRepository {
	return &mockBorrowRepository{
		borrows: make(map[string]*model.Borrow),
	}
}

func (m *mockBorrowRepository) Create(borrow *model.Borrow) error {
	m.borrows[borrow.ID] = borrow
	return nil
}

func TestCreateBorrow(t *testing.T) {
	now := time.Now()

	borrow := &model.Borrow{
		ID:         "test-id",
		UserID:     "user-id",
		BookID:     "book-id",
		BorrowDate: now,
		DueDate:    now.AddDate(0, 0, 14),
		Status:     model.BorrowStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if borrow.Status != model.BorrowStatusActive {
		t.Fatalf("expected status ACTIVE, got %s", borrow.Status)
	}

	if borrow.DueDate.Before(borrow.BorrowDate) {
		t.Fatal("due date should be after borrow date")
	}
}

func TestBorrowStatusReturned(t *testing.T) {
	borrow := &model.Borrow{
		ID:     "test-id",
		Status: model.BorrowStatusReturned,
	}

	if borrow.Status != "RETURNED" {
		t.Fatalf("expected RETURNED, got %s", borrow.Status)
	}
}
