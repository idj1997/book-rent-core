package repo_mocks

import (
	"github.com/idj1997/book-rent-core/domain"
	"github.com/stretchr/testify/mock"
)

type MockedBookRepository struct {
	mock.Mock
}

func (m *MockedBookRepository) GetByID(id int) (*domain.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockedBookRepository) GetByTitle(title string) ([]domain.Book, error) {
	args := m.Called(title)
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (m *MockedBookRepository) Create(book *domain.Book) (uint, error) {
	args := m.Called(book)
	return args.Get(0).(uint), args.Error(1)
}

func (m *MockedBookRepository) Update(book *domain.Book, updates map[string]interface{}) error {
	args := m.Called(book, updates)
	return args.Error(0)
}

func (m *MockedBookRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
