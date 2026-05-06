package model

import(
	"time"
)

type Book struct{
	ID			string `json:"id"`
	Title		string `json:"title"`
	Author		string `json:"author"`
	Category	string `json:"category"`
	Available	bool `json:""available"`
	CreatedAt	time.Time `json:"created_at"`
}

type CreateBookRequest struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Category  string `json:"category"`
	Available bool   `json:"available"`
}

type UpdateBookRequest struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Category  string `json:"category"`
	Available bool   `json:"available"`
}