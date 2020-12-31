package test

import (
	"github.com/idj1997/book-rent-core/config"
	"github.com/idj1997/book-rent-core/domain"
	"github.com/idj1997/book-rent-core/repository"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type BookRepoIntegrationTestSuite struct {
	suite.Suite
	Repo *repository.GormBookRepository
	Db   *gorm.DB
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, &BookRepoIntegrationTestSuite{})
}

func (suite *BookRepoIntegrationTestSuite) SetupSuite() {
	config.InitConfig("test", "../config.yml")
	suite.Db = config.OpenPostgresDB()
}

func (suite *BookRepoIntegrationTestSuite) SetupTest() {
	// start transaction
	tx := suite.Db.Begin()
	// repository
	suite.Repo = repository.NewGormBookRepository(tx)
}

func (suite *BookRepoIntegrationTestSuite) TearDownTest() {
	// rollback current tx
	suite.Repo.Db.Rollback()
	suite.Repo.Db = nil
}

func (suite *BookRepoIntegrationTestSuite) TearDownSuite() {
	config.ClosePostgresDB(suite.Db)
}

func (suite *BookRepoIntegrationTestSuite) TestGetByID_WithInvalidID_ExpectNotFound() {
	a := assert.New(suite.T())
	const ID uint = 5000

	_, err := suite.Repo.GetByID(int(ID))
	a.Error(err)
	a.Equal(domain.NotFound, err.(*domain.RepoError).Type)
}

func (suite *BookRepoIntegrationTestSuite) TestGetByID_WithValidID_ExpectOK() {
	a := assert.New(suite.T())
	const ID uint = 10000

	book, err := suite.Repo.GetByID(int(ID))
	a.Nil(err)
	a.NotNil(book)
	a.Equal(ID, book.ID)
}

func (suite *BookRepoIntegrationTestSuite) TestGetByTitle_WithEmptyTitle_ExpectAll() {
	a := assert.New(suite.T())
	var allBooks []domain.Book

	err := suite.Db.Find(&allBooks).Error
	if err != nil {
		a.FailNow("Error while reading all books: %v\n", err)
	} else {
		books, err := suite.Repo.GetByTitle("")
		a.Equal(len(allBooks), len(books))
		a.Nil(err)
	}
}

func (suite *BookRepoIntegrationTestSuite) TestGetByTitle_WithInvalidTitle_ExpectEmpty() {
	a := assert.New(suite.T())

	books, err := suite.Repo.GetByTitle("invalid title")
	a.Equal(len(books), 0)
	a.Nil(err)
}

func (suite *BookRepoIntegrationTestSuite) TestGetByTitle_WithExactTitle_ExpectOne() {
	a := assert.New(suite.T())
	const title string = "title1"

	books, err := suite.Repo.GetByTitle(title)
	a.Equal(len(books), 1)
	a.Equal(books[0].Title, title)
	a.Nil(err)
}

func (suite *BookRepoIntegrationTestSuite) TestGetByTitle_WithCommonTitle_ExpectMany() {
	a := assert.New(suite.T())
	const title string = "title"

	books, err := suite.Repo.GetByTitle(title)
	a.True(len(books) > 0)
	a.Nil(err)
}

func (suite *BookRepoIntegrationTestSuite) TestCreate_WithValidObject_ExpectOK() {
	a := assert.New(suite.T())
	book := domain.Book{
		Title:   "test title",
		Content: "test content",
		Stock:   10}

	id, err := suite.Repo.Create(&book)
	a.True(id > 0)
	a.Nil(err)
}

func (suite *BookRepoIntegrationTestSuite) TestCreate_WithInvalidObj_ExpectAlreadyExists() {
	a := assert.New(suite.T())
	book, _ := suite.Repo.GetByID(10000)

	_, err := suite.Repo.Create(book)
	a.Error(err)
	a.Equal(domain.UniqueConstraint, err.(*domain.RepoError).Type)
}

func (suite *BookRepoIntegrationTestSuite) TestUpdate_WithEmptyUpdates_ExpectNoChanges() {
	a := assert.New(suite.T())
	originalBook, _ := suite.Repo.GetByID(10000)
	book := *originalBook // shallow copy
	updates := make(map[string]interface{})

	err := suite.Repo.Update(&book, updates)
	a.Nil(err)
	a.Equal(originalBook.ID, book.ID)
	a.Equal(originalBook.Title, book.Title)
	a.Equal(originalBook.Content, book.Content)
	a.Equal(originalBook.Stock, book.Stock)
	a.Equal(originalBook.CreatedAt, book.CreatedAt)
	a.NotEqual(originalBook.UpdatedAt, book.UpdatedAt)
}

func (suite *BookRepoIntegrationTestSuite) TestUpdate_WithStockUpdates_ExpectStockChanged() {
	a := assert.New(suite.T())
	originalBook, _ := suite.Repo.GetByID(10000)
	book := *originalBook // shallow copy
	updates := make(map[string]interface{})
	updates["Stock"] = 0

	err := suite.Repo.Update(&book, updates)
	a.Nil(err)
	a.Equal(originalBook.ID, book.ID)
	a.Equal(originalBook.Title, book.Title)
	a.Equal(originalBook.Content, book.Content)
	a.Equal(originalBook.CreatedAt, book.CreatedAt)
	a.NotEqual(originalBook.UpdatedAt, book.UpdatedAt)
	// main comparison
	a.NotEqual(originalBook.Stock, book.Stock)
	a.Equal(updates["Stock"], book.Stock)
}

func (suite *BookRepoIntegrationTestSuite) TestUpdate_WithInvalidUpdates_ExpectInvalidField() {
	a := assert.New(suite.T())
	originalBook, _ := suite.Repo.GetByID(10000)
	// shallow copy
	book := *originalBook
	updates := make(map[string]interface{})
	updates["InvalidField"] = 0

	err := suite.Repo.Update(&book, updates)
	a.Error(err)
	a.Equal(domain.InvalidField, err.(*domain.RepoError).Type)
}

func (suite *BookRepoIntegrationTestSuite) TestDelete_WithValidID_ExpectDeleted() {
	a := assert.New(suite.T())
	const ID int = 10000

	// delete existing
	err := suite.Repo.Delete(ID)
	a.Nil(err)

	// check if deleted
	_, err = suite.Repo.GetByID(ID)
	a.NotNil(err)
	a.Equal(domain.NotFound, err.(*domain.RepoError).Type)
}
