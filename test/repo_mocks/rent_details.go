package repo_mocks

import (
	"github.com/idj1997/book-rent-core/domain"
	"github.com/stretchr/testify/mock"
)

type MockedRentDetailsRepository struct {
	mock.Mock
}

func (m *MockedRentDetailsRepository) GetByID(id int) (*domain.RentDetails, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.RentDetails), args.Error(1)
}

func (m *MockedRentDetailsRepository) Create(rent *domain.RentDetails) error {
	args := m.Called(rent)
	return args.Error(0)
}

func (m *MockedRentDetailsRepository) Update(rent *domain.RentDetails, updates map[string]interface{}) error {
	args := m.Called(rent, updates)
	return args.Error(0)
}

func (m *MockedRentDetailsRepository) UpdateAssociations(rent *domain.RentDetails, updates map[string]interface{}) error {
	args := m.Called(rent, updates)
	return args.Error(0)
}

func (m *MockedRentDetailsRepository) GetByUser(userID int) ([]domain.RentDetails, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.RentDetails), args.Error(1)
}

func (m *MockedRentDetailsRepository) GetByBook(bookID int) ([]domain.RentDetails, error) {
	args := m.Called(bookID)
	return args.Get(0).([]domain.RentDetails), args.Error(1)
}

func (m *MockedRentDetailsRepository) GetByStatus(status domain.RentDetailsStatus) ([]domain.RentDetails, error) {
	args := m.Called(status)
	return args.Get(0).([]domain.RentDetails), args.Error(1)
}

func (m *MockedRentDetailsRepository) RentDetailsIterator(stream chan domain.RentDetails) {
	defer close(stream)

	rents := make([]domain.RentDetails, 3)
	rents[0].Status = 0
	rents[1].Status = 1
	rents[2].Status = 2

	for _, rent := range rents {
		stream <- rent
	}
}
