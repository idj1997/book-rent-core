package test

import (
	"github.com/idj1997/book-rent-core/domain"
	"github.com/idj1997/book-rent-core/service"
	"github.com/idj1997/book-rent-core/test/repo_mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserServiceUnitTestSuite struct {
	suite.Suite
	service *service.UserService
	repo    *repo_mocks.MockedUserRepository
}

func TestUserServiceUnitTestSuite(t *testing.T) {
	suite.Run(t, &UserServiceUnitTestSuite{})
}

func (suite *UserServiceUnitTestSuite) SetupTest() {
	suite.repo = &repo_mocks.MockedUserRepository{}
	suite.service = &service.UserService{Repo: suite.repo}
}

func (suite *UserServiceUnitTestSuite) TestGetByID_WithInvalidID_ExpectNotFound() {
	a := assert.New(suite.T())
	const ID int = 123124

	suite.repo.
		On("GetByID", ID).
		Return(domain.NilUserPtr, &domain.RepoError{Type: domain.NotFound})

	_, err := suite.service.GetByID(ID)
	a.Error(err)
	a.Equal(service.NotFound, err.(*service.ServiceError).Type)
}

func (suite *UserServiceUnitTestSuite) TestGetByID_WithValidID_ExpectFound() {
	a := assert.New(suite.T())
	const ID int = 10000
	user := domain.User{
		Firstname: "test",
		Lastname:  "test",
		Email:     "test",
		Password:  "test",
		Type:      0,
	}

	suite.repo.
		On("GetByID", ID).
		Return(&user, domain.NilRepoErrPtr)

	returnedUser, err := suite.service.GetByID(ID)
	a.Nil(err)
	a.NotNil(user)
	a.Equal(user.Email, returnedUser.Email)
}

func (suite *UserServiceUnitTestSuite) TestGetByEmail_WithInvalidEmail_ExpectNotFound() {
	a := assert.New(suite.T())
	email := "invalid@email.com"

	suite.repo.
		On("GetByEmail", email).
		Return(domain.NilUserPtr, &domain.RepoError{Type: domain.NotFound})

	_, err := suite.service.GetByEmail(email)
	a.Error(err)
	a.Equal(service.NotFound, err.(*service.ServiceError).Type)
}

func (suite *UserServiceUnitTestSuite) TestGetByEmail_WithValidEmail_ExpectFound() {
	a := assert.New(suite.T())
	user := domain.User{
		Firstname: "test",
		Lastname:  "test",
		Email:     "test",
		Password:  "test",
		Type:      0,
	}

	suite.repo.
		On("GetByEmail", user.Email).
		Return(&user, domain.NilRepoErrPtr)

	returnedUser, err := suite.service.GetByEmail(user.Email)
	a.Nil(err)
	a.NotNil(returnedUser)
	a.Equal(user.Email, returnedUser.Firstname)
}

func (suite *UserServiceUnitTestSuite) TestCreate_WithUnavailableEmail_ExpectAlreadyExist() {
	a := assert.New(suite.T())
	user := domain.User{
		Firstname: "test",
		Lastname:  "test",
		Email:     "unavailable@gmail.com",
		Password:  "test",
		Type:      0,
	}
	repoError := domain.RepoError{Type: domain.UniqueConstraint}

	suite.repo.
		On("Create", &user).
		Return(&repoError)

	serviceErr := suite.service.Create(&user)
	a.Error(serviceErr)
	a.Equal(service.AlreadyExist, serviceErr.(*service.ServiceError).Type)
}

func (suite *UserServiceUnitTestSuite) TestCreate_WithUnavailableID_ExpectAlreadyExist() {
	a := assert.New(suite.T())
	user := domain.User{
		Model:     gorm.Model{ID: 10000},
		Firstname: "test",
		Lastname:  "test",
		Email:     "available@gmail.com",
		Password:  "test",
		Type:      0,
	}
	repoError := domain.RepoError{Type: domain.UniqueConstraint}

	suite.repo.
		On("Create", &user).
		Return(&repoError)

	serviceErr := suite.service.Create(&user)
	a.Error(serviceErr)
	a.Equal(service.AlreadyExist, serviceErr.(*service.ServiceError).Type)
}

func (suite *UserServiceUnitTestSuite) TestCreate_WithValidUserObj_ExpectCreated() {
	a := assert.New(suite.T())
	user := domain.User{
		Firstname: "test",
		Lastname:  "test",
		Email:     "available@gmail.com",
		Password:  "test",
		Type:      0,
	}

	suite.repo.
		On("Create", &user).
		Return(domain.NilRepoErrPtr)

	err := suite.service.Create(&user)
	a.Nil(err)
}

func (suite *UserServiceUnitTestSuite) TestDelete_WithInvalidID_ExpectNotFound() {
	a := assert.New(suite.T())
	id := 11111
	repoErr := domain.RepoError{Type: domain.NotFound}

	suite.repo.
		On("GetByID", id).
		Return(domain.NilUserPtr, &repoErr)

	serviceErr := suite.service.Delete(id)
	a.Error(serviceErr)
	a.Equal(service.NotFound, serviceErr.(*service.ServiceError).Type)
}

func (suite *UserServiceUnitTestSuite) TestDelete_WithValidID_ExpectDeleted() {
	a := assert.New(suite.T())
	id := 1
	userPtr := &domain.User{
		Model: gorm.Model{
			ID:        1,
			CreatedAt: time.Time{},
			UpdatedAt: time.Time{},
			DeletedAt: gorm.DeletedAt{},
		},
		Firstname: "test",
		Lastname:  "test",
		Email:     "test",
		Password:  "test",
		Type:      0,
	}

	suite.repo.
		On("GetByID", id).
		Return(userPtr, domain.NilRepoErrPtr)

	suite.repo.
		On("Delete", id).
		Return(domain.NilRepoErrPtr)

	serviceErr := suite.service.Delete(id)
	a.Nil(serviceErr)
}
