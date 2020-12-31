package service

import (
	"github.com/idj1997/book-rent-core/domain"
)

type UserService struct {
	Repo domain.UserRepository
}

func (u *UserService) GetByID(id int) (*domain.User, error) {
	user, err := u.Repo.GetByID(id)
	return user, RepoErrorToServiceError(err)
}

func (u *UserService) GetByEmail(email string) (*domain.User, error) {
	user, err := u.Repo.GetByEmail(email)
	return user, RepoErrorToServiceError(err)
}

func (u *UserService) GetByFirstnameAndLastname(firstname string, lastname string) ([]domain.User, error) {
	panic("implement me")
}

func (u *UserService) Create(user *domain.User) error {
	err := u.Repo.Create(user)
	return RepoErrorToServiceError(err)
}

func (u *UserService) Delete(id int) error {
	_, err := u.GetByID(id)
	if err != nil {
		return err
	}

	err = u.Repo.Delete(id)
	return RepoErrorToServiceError(err)
}
