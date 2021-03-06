package service

import (
	"github.com/go-playground/validator"
	"github.com/idj1997/book-rent-core/domain"
)

type BookService struct {
	br domain.BookRepository
}

func NewBookService(br domain.BookRepository) *BookService {
	return &BookService{br: br}
}

func (bs *BookService) GetByID(id int) (*domain.Book, error) {
	book, err := bs.br.GetByID(id)
	return book, RepoErrorToServiceError(err)
}

func (bs *BookService) GetByTitle(title string) ([]domain.Book, error) {
	books, err := bs.br.GetByTitle(title)
	return books, RepoErrorToServiceError(err)
}

func (bs *BookService) Create(book *domain.Book) (int, error) {
	validate := validator.New()
	validationErr := validate.Struct(book)
	if validationErr != nil {
		return 0, &ServiceError{Type: InvalidArguments}
	}

	id, err := bs.br.Create(book)
	return int(id), RepoErrorToServiceError(err)
}

func (bs *BookService) UpdateStock(bookID int, newStock int) (*domain.Book, error) {
	if newStock <= 0 {
		return nil, &ServiceError{Type: InvalidArguments}
	}

	book, err := bs.GetByID(bookID)
	if err != nil {
		return nil, err
	}

	stockUpdates := make(map[string]interface{})
	stockUpdates["Stock"] = newStock
	err = bs.br.Update(book, stockUpdates)
	return book, RepoErrorToServiceError(err)
}

func (bs *BookService) Delete(id int) error {
	_, err := bs.br.GetByID(id)
	if err != domain.NilRepoErrPtr {
		return RepoErrorToServiceError(err)
	}

	err = bs.br.Delete(id)
	return RepoErrorToServiceError(err)
}
