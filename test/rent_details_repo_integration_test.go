package test

import (
	"github.com/idj1997/book-rent-core/config"
	"github.com/idj1997/book-rent-core/domain"
	"github.com/idj1997/book-rent-core/repository"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RentDetailsIntegrationTestSuite struct {
	suite.Suite
	Repo *repository.GormRentDetailsRepository
	Db   *gorm.DB
}

func TestRentDetailsIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &RentDetailsIntegrationTestSuite{})
}

func (suite *RentDetailsIntegrationTestSuite) SetupSuite() {
	config.InitConfig("test", "../config.yml")
	suite.Db = config.OpenPostgresDB()
}

func (suite *RentDetailsIntegrationTestSuite) SetupTest() {
	tx := suite.Db.Begin()
	suite.Repo = &repository.GormRentDetailsRepository{Db: tx}
}

func (suite *RentDetailsIntegrationTestSuite) TearDownTest() {
	suite.Repo.Db.Rollback()
	suite.Repo.Db = nil
}

func (suite *RentDetailsIntegrationTestSuite) TearDownSuite() {
	config.ClosePostgresDB(suite.Db)
}

func (suite *RentDetailsIntegrationTestSuite) TestGetByID_WithValidID_ExpectOk() {
	a := assert.New(suite.T())
	id := 10000

	rent, err := suite.Repo.GetByID(id)
	a.Nil(err)
	a.Equal(uint(id), rent.ID)
	a.Equal(uint(10000), rent.User.ID)
	a.Equal(uint(10000), rent.Book.ID)
	a.Equal("john", rent.User.Firstname)
	a.Equal("title1", rent.Book.Title)
}

func (suite *RentDetailsIntegrationTestSuite) TestCreate_WithValidObj_ExpectOk() {
	a := assert.New(suite.T())
	rent := domain.RentDetails{
		UserID: 10000, // user: john doe
		BookID: 10000, // book: title 1
		Status: domain.RENTED,
	}

	repoErr := suite.Repo.Create(&rent)
	a.Nil(repoErr)
	a.NotNil(rent.ID)

	fullRent, _ := suite.Repo.GetByID(int(rent.ID))
	a.Equal("title1", fullRent.Book.Title)
	a.Equal("john", fullRent.User.Firstname)
}

func (suite *RentDetailsIntegrationTestSuite) TestCreate_WithInvalidObj_ExpectForeignKeyConstraint() {
	a := assert.New(suite.T())
	rent := domain.RentDetails{
		UserID: 0,
		BookID: 0,
		Status: domain.RENTED,
	}

	repoErr := suite.Repo.Create(&rent)
	a.Error(repoErr)
	a.Equal(domain.ForeignKeyConstraint, repoErr.(*domain.RepoError).Type)
}

func (suite *RentDetailsIntegrationTestSuite) TestUpdate_WithInvalidForeignKey_ExpectForeignKeyConstraint() {
	a := assert.New(suite.T())
	id := 10000
	rent, _ := suite.Repo.GetByID(id)
	updates := make(map[string]interface{})
	updates["UserID"] = 0

	repoErr := suite.Repo.Update(rent, updates)
	a.Error(repoErr)
	a.Equal(domain.ForeignKeyConstraint, repoErr.(*domain.RepoError).Type)
}

func (suite *RentDetailsIntegrationTestSuite) TestUpdate_WithInvalidField_ExpectInvalidFieldErr() {
	a := assert.New(suite.T())
	id := 10000
	rent, _ := suite.Repo.GetByID(id)
	updates := make(map[string]interface{})
	updates["InvalidField"] = 0

	repoErr := suite.Repo.Update(rent, updates)
	a.Error(repoErr)
	a.Equal(domain.InvalidField, repoErr.(*domain.RepoError).Type)
}

func (suite *RentDetailsIntegrationTestSuite) TestUpdate_WithValidUpdates_ExpectOk() {
	a := assert.New(suite.T())
	rentId := 10000
	rent, _ := suite.Repo.GetByID(rentId)

	newStatus := domain.RETURNED
	updates := make(map[string]interface{})
	updates["Status"] = newStatus

	repoErr := suite.Repo.Update(rent, updates)
	a.Nil(repoErr)
	a.Equal(newStatus, rent.Status)
}

func (suite *RentDetailsIntegrationTestSuite) TestUpdate_WithValidAssocUpdates_ExpectOk() {
	a := assert.New(suite.T())
	rentId := 10000
	rent, _ := suite.Repo.GetByID(rentId)

	newBookID := 10001
	updates := make(map[string]interface{})
	updates["book_id"] = newBookID

	repoErr := suite.Repo.UpdateAssociations(rent, updates)
	a.Nil(repoErr)
	a.Equal(uint(newBookID), rent.Book.ID)
}

func (suite *RentDetailsIntegrationTestSuite) TestGetByUser_WithInvalidID_ExpectEmpty() {
	a := assert.New(suite.T())
	userID := 31213

	rents, err := suite.Repo.GetByUser(userID)
	a.Nil(err)
	a.NotNil(rents)
	a.Empty(rents)
}

func (suite *RentDetailsIntegrationTestSuite) TestGetByUser_WithValidID_ExpectMany() {
	a := assert.New(suite.T())
	userID := 10000

	rents, err := suite.Repo.GetByUser(userID)
	a.Nil(err)
	a.NotNil(rents)
	a.NotEmpty(rents)
}

func (suite *RentDetailsIntegrationTestSuite) TestGetByBook_WithInvalidID_ExpectEmpty() {
	a := assert.New(suite.T())
	bookID := 31213

	rents, err := suite.Repo.GetByBook(bookID)
	a.Nil(err)
	a.NotNil(rents)
	a.Empty(rents)
}

func (suite *RentDetailsIntegrationTestSuite) TestGetByBook_WithValidID_ExpectMany() {
	a := assert.New(suite.T())
	bookID := 10000

	rents, err := suite.Repo.GetByBook(bookID)
	a.Nil(err)
	a.NotNil(rents)
	a.NotEmpty(rents)
}

func (suite *RentDetailsIntegrationTestSuite) TestGetByStatus_WithValidStatus_ExpectMany() {
	a := assert.New(suite.T())
	status := domain.RENTED

	rents, err := suite.Repo.GetByStatus(status)
	a.Nil(err)
	a.NotNil(rents)
	a.NotEmpty(rents)
}

func (suite *RentDetailsIntegrationTestSuite) TestGetByStatus_WithInvalidStatus_ExpectEmpty() {
	a := assert.New(suite.T())
	status := 5124

	rents, err := suite.Repo.GetByStatus(domain.RentDetailsStatus(status))
	a.Nil(err)
	a.Empty(rents)
}

func (suite *RentDetailsIntegrationTestSuite) TestRentedAndExpiredProducer_ExpectMany() {
	a := assert.New(suite.T())
	rents := make([]domain.RentDetails, 0)
	stream := make(chan domain.RentDetails)

	go suite.Repo.RentDetailsIterator(stream)
	for rent := range stream {
		rents = append(rents, rent)
	}

	a.NotEmpty(rents)
}
