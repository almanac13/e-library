package model

import "time"

type Borrow struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	BookID     string     `json:"book_id"`
	BorrowDate time.Time  `json:"borrow_date"`
	DueDate    time.Time  `json:"due_date"`
	ReturnDate *time.Time `json:"return_date,omitempty"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

const (
	BorrowStatusActive   = "ACTIVE"
	BorrowStatusReturned = "RETURNED"
)
