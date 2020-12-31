package test

import (
	"github.com/idj1997/book-rent-core/domain"
	"github.com/idj1997/book-rent-core/service"
	"github.com/idj1997/book-rent-core/test/repo_mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BookServiceUnitTestSuite struct {
	suite.Suite
	service *service.BookService
	repo    *repo_mocks.MockedBookRepository
}

func TestBookServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, &BookServiceUnitTestSuite{})
}

func (suite *BookServiceUnitTestSuite) SetupTest() {
	suite.repo = &repo_mocks.MockedBookRepository{}
	suite.service = service.NewBookService(suite.repo)
}

func (suite *BookServiceUnitTestSuite) TestGetByID_WithInvalidId_ExpectNotFound() {
	a := assert.New(suite.T())
	invalidID := -1
	var bookPtr *domain.Book = nil

	suite.repo.On("GetByID", invalidID).Return(bookPtr, &domain.RepoError{Type: domain.NotFound})

	book, err := suite.service.GetByID(invalidID)
	a.Nil(book)
	a.Error(err)
	a.Equal(service.NotFound, err.(*service.ServiceError).Type)
}

func (suite *BookServiceUnitTestSuite) TestGetByTitle_WithInvalidTitle_ExpectEmpty() {
	a := assert.New(suite.T())
	title := "invalid title"
	empty := make([]domain.Book, 0)

	suite.repo.
		On("GetByTitle", title).
		Return(empty, domain.NilRepoErrPtr)

	books, err := suite.service.GetByTitle(title)
	a.Nil(err)
	a.Empty(books)
}

func (suite *BookServiceUnitTestSuite) TestGetByID_WithValidId_ExpectOk() {
	a := assert.New(suite.T())
	ID := 1
	book := domain.Book{Title: "generic title", Content: "generic test"}

	suite.repo.
		On("GetByID", ID).
		Return(&book, domain.NilRepoErrPtr)

	resultBook, err := suite.service.GetByID(ID)
	a.NotNil(resultBook)
	a.Nil(err)
	a.Equal(book.Title, resultBook.Title)
	a.Equal(book.Content, resultBook.Content)
}

func (suite *BookServiceUnitTestSuite) TestCreate_WithValidObj_ExpectOk() {
	a := assert.New(suite.T())
	book := domain.Book{
		Title:   "test title",
		Content: "test content",
		Stock:   10}

	shouldCreateBook := book
	shouldCreateBook.ID = 1
	suite.repo.
		On("Create", &book).
		Return(shouldCreateBook.ID, domain.NilRepoErrPtr)

	createdBookID, err := suite.service.Create(&book)
	a.Nil(err)
	a.Equal(int(shouldCreateBook.ID), createdBookID)
}

func (suite *BookServiceUnitTestSuite) TestCreate_WithUnavailableID_ExpectAlreadyExists() {
	a := assert.New(suite.T())
	book := domain.Book{
		Title:   "test title",
		Content: "test content",
		Stock:   10}
	book.ID = 2

	suite.repo.
		On("Create", &book).
		Return(uint(0), &domain.RepoError{Type: domain.UniqueConstraint})

	createdBookID, err := suite.service.Create(&book)
	a.Zero(createdBookID)
	a.NotNil(err)
	a.Equal(service.AlreadyExist, err.(*service.ServiceError).Type)
}

func (suite *BookServiceUnitTestSuite) TestUpdateStock_WithInvalidStockValue_ExpectInvalidStockCount() {
	a := assert.New(suite.T())
	book := domain.Book{
		Title:   "test title",
		Content: "test content",
		Stock:   10}
	book.ID = 1
	invalidStockCount := -1

	_, err := suite.service.UpdateStock(int(book.ID), invalidStockCount)
	a.NotNil(err)
	a.Equal(service.InvalidArguments, err.(*service.ServiceError).Type)

}

func (suite *BookServiceUnitTestSuite) TestUpdateStock_WithInvalidBookID_ExpectNotFound() {
	a := assert.New(suite.T())
	book := domain.Book{
		Title:   "test title",
		Content: "test content",
		Stock:   10}
	book.ID = 12134
	var bookPtr *domain.Book = nil

	suite.repo.
		On("GetByID", int(book.ID)).
		Return(bookPtr, &domain.RepoError{Type: domain.NotFound})

	_, err := suite.service.UpdateStock(int(book.ID), 120)
	a.NotNil(err)
	a.Equal(service.NotFound, err.(*service.ServiceError).Type)
}

func (suite *BookServiceUnitTestSuite) TestUpdateStock_WithValidArgs_ExpectOk() {
	a := assert.New(suite.T())
	book := domain.Book{
		Title:   "test title",
		Content: "test content",
		Stock:   10}
	book.ID = 1
	newStockCount := 15

	suite.repo.
		On("GetByID", int(book.ID)).
		Return(&book, domain.NilRepoErrPtr)

	updates := make(map[string]interface{})
	updates["Stock"] = newStockCount

	suite.repo.
		On("Update", &book, updates).
		Return(domain.NilRepoErrPtr)

	_, err := suite.service.UpdateStock(int(book.ID), newStockCount)
	a.Nil(err)
}
