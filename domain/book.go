package domain

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title   string
	Content string
	Stock   int
}

type BookRepository interface {
	GetByID(id int) (*Book, error)
	GetByTitle(title string) ([]Book, error)
	Create(book *Book) (uint, error)
	Update(book *Book, updates map[string]interface{}) error
	Delete(id int) error
}

type BookService interface {
	GetByID(id int) (*Book, error)
	GetByTitle(title string) ([]Book, error)
	Create(book *Book) (int, error)
	UpdateStock(bookID int, newStock int) (*Book, error)
	Delete(id int) error
}
