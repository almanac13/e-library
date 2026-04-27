package model

type Book struct{
	ID			string `json:"id"`
	Title		string `json:"title"`
	Author		string `json:"author"`
	Category	string `json:"category"`
	Available	bool `json:""available`
}