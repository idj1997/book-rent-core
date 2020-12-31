package repo_mocks

import (
	"github.com/idj1997/book-rent-core/domain"
	"github.com/stretchr/testify/mock"
)

type MockedUserRepository struct {
	mock.Mock
}

func (m *MockedUserRepository) GetByID(id int) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockedUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)

}

func (m *MockedUserRepository) GetByFirstnameAndLastname(firstname string, lastname string) ([]domain.User, error) {
	args := m.Called(firstname, lastname)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockedUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockedUserRepository) Update(user *domain.User, updates map[string]interface{}) error {
	args := m.Called(user, updates)
	return args.Error(0)
}

func (m *MockedUserRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
