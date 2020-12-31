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

type UserRepoIntegrationTestSuite struct {
	suite.Suite
	Repo *repository.GormUserRepository
	Db   *gorm.DB
}

func TestRunUserRepoIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &UserRepoIntegrationTestSuite{})
}

func (suite *UserRepoIntegrationTestSuite) SetupSuite() {
	config.InitConfig("test", "../config.yml")
	suite.Db = config.OpenPostgresDB()
}

func (suite *UserRepoIntegrationTestSuite) SetupTest() {
	// start transaction
	tx := suite.Db.Begin()
	// repository
	suite.Repo = repository.NewGormUserRepository(tx)
}

func (suite *UserRepoIntegrationTestSuite) TearDownTest() {
	// rollback current tx
	suite.Repo.Db.Rollback()
	suite.Repo.Db = nil
}

func (suite *UserRepoIntegrationTestSuite) TearDownSuite() {
	config.ClosePostgresDB(suite.Db)
}

func (suite *UserRepoIntegrationTestSuite) TestGetByID_WithInvalidID_ExpectNotFound() {
	a := assert.New(suite.T())
	const ID uint = 5000

	_, err := suite.Repo.GetByID(int(ID))
	a.Error(err)
	a.Equal(domain.NotFound, err.(*domain.RepoError).Type)
}

func (suite *UserRepoIntegrationTestSuite) TestGetByID_WithValidID_ExpectOK() {
	a := assert.New(suite.T())
	const ID uint = 10000

	user, err := suite.Repo.GetByID(int(ID))
	a.Nil(err)
	a.NotNil(user)
	a.Equal(ID, user.ID)
}

func (suite *UserRepoIntegrationTestSuite) TestGetByID_WithInvalidEmail_ExpectNotFound() {
	a := assert.New(suite.T())
	email := "invalid@gmail.com"

	_, err := suite.Repo.GetByEmail(email)
	a.Error(err)
	a.Equal(domain.NotFound, err.(*domain.RepoError).Type)
}

func (suite *UserRepoIntegrationTestSuite) TestGetByID_WithValidEmail_ExpectOK() {
	a := assert.New(suite.T())
	email := "johndoe@gmail.com"

	user, err := suite.Repo.GetByEmail(email)
	a.Nil(err)
	a.NotNil(user)
	a.Equal(email, user.Email)
}

func (suite *UserRepoIntegrationTestSuite) TestGetByFirstnameAndLastname_WithEmptyArgs_ExpectAll() {
	a := assert.New(suite.T())
	var allUsers []domain.User

	err := suite.Db.Find(&allUsers).Error
	if err != nil {
		a.FailNow("Error while reading all users: %v\n", err)
	} else {
		books, err := suite.Repo.GetByFirstnameAndLastname("", "")
		a.Equal(len(allUsers), len(books))
		a.Nil(err)
	}
}

func (suite *UserRepoIntegrationTestSuite) TestGetByFirstnameAndLastname_WithInvalidArgs_ExpectEmpty() {
	a := assert.New(suite.T())

	firstname := "invalid firstname"
	lastname := "invalid lastname"
	users, err := suite.Repo.GetByFirstnameAndLastname(firstname, lastname)
	a.Equal(len(users), 0)
	a.Nil(err)
}

func (suite *UserRepoIntegrationTestSuite) TestGetByFirstnameAndLastname_WithExactArgs_ExpectOne() {
	a := assert.New(suite.T())
	firstname := "john"
	lastname := "doe"

	users, err := suite.Repo.GetByFirstnameAndLastname(firstname, lastname)
	a.Equal(len(users), 1)
	a.Equal(users[0].Firstname, firstname)
	a.Equal(users[0].Lastname, lastname)
	a.Nil(err)
}

func (suite *UserRepoIntegrationTestSuite) TestCreate_WithValidObject_ExpectOK() {
	a := assert.New(suite.T())
	user := domain.User{
		Firstname: "test",
		Lastname:  "test",
		Email:     "test@test.com",
		Password:  "$2y$12$Z51tvYyB2xEUejQydGcaiuCs1i3xqgHMvHwlzVLLQCk/7KVzahP9W",
		Type:      domain.CUSTOMER}

	err := suite.Repo.Create(&user)
	a.True(user.ID > 0)
	a.Equal("test", user.Firstname)
	a.Nil(err)
}

func (suite *UserRepoIntegrationTestSuite) TestCreate_WithUnavailableID_ExpectAlreadyExists() {
	a := assert.New(suite.T())
	user, _ := suite.Repo.GetByID(10000)

	err := suite.Repo.Create(user)
	a.Error(err)
	a.Equal(domain.UniqueConstraint, err.(*domain.RepoError).Type)
}

func (suite *UserRepoIntegrationTestSuite) TestCreate_WithUnavailableEmail_ExpectAlreadyExists() {
	a := assert.New(suite.T())
	user := domain.User{
		Firstname: "test",
		Lastname:  "test",
		Email:     "johndoe@gmail.com",
		Password:  "$2y$12$Z51tvYyB2xEUejQydGcaiuCs1i3xqgHMvHwlzVLLQCk/7KVzahP9W",
		Type:      domain.CUSTOMER}

	err := suite.Repo.Create(&user)
	a.Error(err)
	a.Equal(domain.UniqueConstraint, err.(*domain.RepoError).Type)
}

func (suite *UserRepoIntegrationTestSuite) TestUpdate_WithEmptyUpdates_ExpectNoChanges() {
	a := assert.New(suite.T())
	originalUser, _ := suite.Repo.GetByID(10000)
	user := *originalUser // shallow copy
	updates := make(map[string]interface{})

	err := suite.Repo.Update(&user, updates)
	a.Nil(err)
	a.Equal(originalUser.ID, user.ID)
	a.Equal(originalUser.Firstname, user.Firstname)
	a.Equal(originalUser.Lastname, user.Lastname)
	a.Equal(originalUser.Email, user.Email)
	a.Equal(originalUser.Password, user.Password)
	a.Equal(originalUser.CreatedAt, user.CreatedAt)
	a.NotEqual(originalUser.UpdatedAt, user.UpdatedAt)
}

func (suite *UserRepoIntegrationTestSuite) TestUpdate_WithNewPassword_ExpectPasswordChanged() {
	a := assert.New(suite.T())
	originalUser, _ := suite.Repo.GetByID(10000)
	user := *originalUser // shallow copy
	updates := make(map[string]interface{})
	updates["Password"] = "test"

	err := suite.Repo.Update(&user, updates)
	a.Nil(err)
	a.Equal(originalUser.ID, user.ID)
	a.Equal(originalUser.Firstname, user.Firstname)
	a.Equal(originalUser.Lastname, user.Lastname)
	a.Equal(originalUser.Email, user.Email)
	a.NotEqual(originalUser.Password, user.Password)
	a.Equal("test", user.Password)
	a.Equal(originalUser.CreatedAt, user.CreatedAt)
	a.NotEqual(originalUser.UpdatedAt, user.UpdatedAt)
}

func (suite *UserRepoIntegrationTestSuite) TestUpdate_WithInvalidUpdates_ExpectInvalidField() {
	a := assert.New(suite.T())
	originalUser, _ := suite.Repo.GetByID(10000)
	// shallow copy
	user := *originalUser
	updates := make(map[string]interface{})
	updates["InvalidField"] = 0

	err := suite.Repo.Update(&user, updates)
	a.Error(err)
	a.Equal(domain.InvalidField, err.(*domain.RepoError).Type)
}

func (suite *UserRepoIntegrationTestSuite) TestUpdate_WithUnavailableEmail_ExpectUniqueConstraint() {
	a := assert.New(suite.T())
	originalUser, _ := suite.Repo.GetByID(10001)
	// shallow copy
	user := *originalUser
	updates := make(map[string]interface{})
	updates["Email"] = "johndoe@gmail.com"

	err := suite.Repo.Update(&user, updates)
	a.Error(err)
	a.Equal(domain.UniqueConstraint, err.(*domain.RepoError).Type)
}

func (suite *UserRepoIntegrationTestSuite) TestDelete_WithValidID_ExpectDeleted() {
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
