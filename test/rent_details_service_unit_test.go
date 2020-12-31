package test

import (
	"github.com/idj1997/book-rent-core/domain"
	"github.com/idj1997/book-rent-core/service"
	"github.com/idj1997/book-rent-core/test/repo_mocks"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RentDetailsUnitTestSuite struct {
	suite.Suite
	RentService domain.RentDetailsService
	RentRepo    *repo_mocks.MockedRentDetailsRepository
	BookRepo    *repo_mocks.MockedBookRepository
}

func TestRentDetailsUnitTestSuite(t *testing.T) {
	suite.Run(t, &RentDetailsUnitTestSuite{})
}

func (suite *RentDetailsUnitTestSuite) SetupTest() {
	suite.RentRepo = &repo_mocks.MockedRentDetailsRepository{}
	suite.BookRepo = &repo_mocks.MockedBookRepository{}
	suite.RentService = &service.RentDetailsService{
		RentRepo: suite.RentRepo,
		BookRepo: suite.BookRepo}
}

func (suite *RentDetailsUnitTestSuite) TestGetByID_WithInvalidRentID_ExpectNotFound() {
	a := assert.New(suite.T())
	id := 123

	suite.RentRepo.
		On("GetByID", id).
		Return(domain.NilRentPtr, &domain.RepoError{Type: domain.NotFound})

	rent, err := suite.RentService.GetByID(id)
	a.Nil(rent)
	a.Equal(service.NotFound, err.(*service.ServiceError).Type)
}

func (suite *RentDetailsUnitTestSuite) TestGetByID_WithValidRentID_ExpectFound() {
	a := assert.New(suite.T())
	id := 10000 // valid
	rent := domain.RentDetails{
		UserID: 123,
		BookID: 132,
		Status: 1}
	rent.ID = uint(id)

	suite.RentRepo.
		On("GetByID", id).
		Return(&rent, domain.NilRepoErrPtr)

	returnedRent, err := suite.RentService.GetByID(id)
	a.NotNil(returnedRent)
	a.Nil(err)
	a.Equal(rent.UserID, returnedRent.UserID)
}

func (suite *RentDetailsUnitTestSuite) TestRentBook_WithInvalidBookID_ExpectNotFound() {
	a := assert.New(suite.T())

	rent := domain.RentDetails{
		UserID: 10000, // valid id
		BookID: 0,     // invalid id
		Status: 0}

	suite.BookRepo.
		On("GetByID", rent.BookID).
		Return(domain.NilBookPtr, &domain.RepoError{Type: domain.NotFound})

	err := suite.RentService.RentBook(&rent)
	a.Equal(service.NotFound, err.(*service.ServiceError).Type)
}

func (suite *RentDetailsUnitTestSuite) TestRentBook_WithEmptyBookStock_ExpectBookNotAvailable() {
	a := assert.New(suite.T())

	book := domain.Book{
		Title:   "test",
		Content: "test",
		Stock:   0, // empty
	}

	rent := domain.RentDetails{
		UserID: 10000, // valid id
		BookID: 10000, // valid id
		Status: 0}

	suite.BookRepo.
		On("GetByID", rent.BookID).
		Return(&book, domain.NilRepoErrPtr)

	err := suite.RentService.RentBook(&rent)
	a.NotNil(err)
	a.Equal(service.NotEnoughBooksOnStock, err.(*service.ServiceError).Type)
}

func (suite *RentDetailsUnitTestSuite) TestRentBook_WithInvalidUserID_ExpectNotFound() {
	a := assert.New(suite.T())

	book := domain.Book{
		Title:   "test",
		Content: "test",
		Stock:   10, // available
	}

	rent := domain.RentDetails{
		UserID: 0,     // invalid id
		BookID: 10000, // valid id
		Status: 0}

	suite.BookRepo.
		On("GetByID", rent.BookID).
		Return(&book, domain.NilRepoErrPtr)

	suite.RentRepo.
		On("Create", &rent).
		Return(&domain.RepoError{Type: domain.NotFound})

	err := suite.RentService.RentBook(&rent)
	a.Equal(service.NotFound, err.(*service.ServiceError).Type)
}

func (suite *RentDetailsUnitTestSuite) TestRentBook_WithValidObj_ExpectCreated() {
	a := assert.New(suite.T())

	book := domain.Book{
		Title:   "test",
		Content: "test",
		Stock:   10, // available
	}

	rent := domain.RentDetails{
		UserID: 10000, // valid id
		BookID: 10000, // valid id
		Status: 0}

	updates := make(map[string]interface{})
	updates["stock"] = book.Stock - 1

	suite.BookRepo.
		On("GetByID", rent.BookID).
		Return(&book, domain.NilRepoErrPtr)

	suite.RentRepo.
		On("Create", &rent).
		Return(domain.NilRepoErrPtr)

	suite.BookRepo.
		On("Update", &book, updates).
		Return(domain.NilRepoErrPtr)

	err := suite.RentService.RentBook(&rent)
	a.Nil(err)
	a.True(rent.ReturnDeadline.After(rent.CreatedAt))
}

func (suite *RentDetailsUnitTestSuite) TestReturnBook_WithInvalidID_ExpectNotFound() {
	a := assert.New(suite.T())
	id := 312412
	err := domain.RepoError{Type: domain.NotFound}

	suite.RentRepo.
		On("GetByID", id).
		Return(domain.NilRentPtr, &err)

	serviceErr := suite.RentService.ReturnBook(id)
	a.NotNil(serviceErr)
	a.Equal(service.NotFound, serviceErr.(*service.ServiceError).Type)
}

func (suite *RentDetailsUnitTestSuite) TestReturnBook_WithReturnedBook_ExpectInvalid() {
	a := assert.New(suite.T())
	id := 10000
	rent := domain.RentDetails{
		UserID: 100,
		BookID: 100,
		Status: 1}

	suite.RentRepo.
		On("GetByID", id).
		Return(&rent, domain.NilRepoErrPtr)

	err := suite.RentService.ReturnBook(id)
	a.NotNil(err)
	a.Equal(service.BookAlreadyReturned, err.(*service.ServiceError).Type)
}

func (suite *RentDetailsUnitTestSuite) TestReturnBook_WithValidID_ExpectOK() {
	a := assert.New(suite.T())
	id := 10000

	book := domain.Book{
		Title:   "test",
		Content: "test",
		Stock:   10}

	rent := domain.RentDetails{
		UserID: 100,
		BookID: 100,
		Status: 0, // book is not returned
		Book:   book}

	suite.RentRepo.
		On("GetByID", id).
		Return(&rent, domain.NilRepoErrPtr)

	rentUpdates := make(map[string]interface{})
	rentUpdates["status"] = domain.RETURNED
	suite.RentRepo.
		On("Update", &rent, rentUpdates).
		Return(domain.NilRepoErrPtr)

	bookUpdates := make(map[string]interface{})
	bookUpdates["stock"] = book.Stock + 1
	suite.BookRepo.
		On("Update", &book, bookUpdates).
		Return(domain.NilRepoErrPtr)

	err := suite.RentService.ReturnBook(id)
	a.Nil(err)
}

func (suite *RentDetailsUnitTestSuite) TestGetByUser_WithInvalidID_ExpectEmpty() {
	a := assert.New(suite.T())
	id := 1124123
	var rents []domain.RentDetails

	suite.RentRepo.
		On("GetByUser", id).
		Return(rents, domain.NilRepoErrPtr)

	returnedRents, err := suite.RentService.GetByUser(id)
	a.Nil(err)
	a.Empty(returnedRents)
}

func (suite *RentDetailsUnitTestSuite) TestGetByUser_WithValidID_ExpectMany() {
	a := assert.New(suite.T())
	id := 10000
	rents := make([]domain.RentDetails, 2)

	suite.RentRepo.
		On("GetByUser", id).
		Return(rents, domain.NilRepoErrPtr)

	returnedRents, err := suite.RentService.GetByUser(id)
	a.Nil(err)
	a.NotEmpty(returnedRents)
}

func (suite *RentDetailsUnitTestSuite) TestGetByBook_WithInvalidID_ExpectEmpty() {
	a := assert.New(suite.T())
	id := 1124123
	var rents []domain.RentDetails

	suite.RentRepo.
		On("GetByBook", id).
		Return(rents, domain.NilRepoErrPtr)

	returnedRents, err := suite.RentService.GetByBook(id)
	a.Nil(err)
	a.Empty(returnedRents)
}

func (suite *RentDetailsUnitTestSuite) TestGetByBook_WithValidID_ExpectEmpty() {
	a := assert.New(suite.T())
	id := 10000
	rents := make([]domain.RentDetails, 2)

	suite.RentRepo.
		On("GetByBook", id).
		Return(rents, domain.NilRepoErrPtr)

	returnedRents, err := suite.RentService.GetByBook(id)
	a.Nil(err)
	a.NotEmpty(returnedRents)
}

func (suite *RentDetailsUnitTestSuite) TestGetByStatus_WithInvalidStatus_ExpectEmpty() {
	a := assert.New(suite.T())
	status := domain.RentDetailsStatus(1124123)
	var rents []domain.RentDetails

	suite.RentRepo.
		On("GetByStatus", status).
		Return(rents, domain.NilRepoErrPtr)

	returnedRents, err := suite.RentService.GetByStatus(status)
	a.Nil(err)
	a.Empty(returnedRents)
}

func (suite *RentDetailsUnitTestSuite) TestGetByStatus_WithValidID_ExpectEmpty() {
	a := assert.New(suite.T())
	status := domain.RETURNED
	rents := make([]domain.RentDetails, 2)

	suite.RentRepo.
		On("GetByStatus", status).
		Return(rents, domain.NilRepoErrPtr)

	returnedRents, err := suite.RentService.GetByStatus(status)
	a.Nil(err)
	a.NotEmpty(returnedRents)
}

func (suite *RentDetailsUnitTestSuite) TestUpdateToExpired_ExpectMany() {
	a := assert.New(suite.T())

	suite.RentRepo.
		On("Update", mock.Anything, mock.Anything).
		Return(domain.NilRepoErrPtr)

	err := suite.RentService.UpdateToExpired()
	a.Nil(err)
}
