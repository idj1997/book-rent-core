package repository

import (
	"github.com/idj1997/book-rent-core/domain"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	Db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{Db: db}
}

func (repo *GormUserRepository) GetByID(id int) (*domain.User, error) {
	var user domain.User
	err := repo.Db.First(&user, id).Error
	return &user, ErrorToRepoError(err)
}

func (repo *GormUserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := repo.Db.Where("email = ?", email).First(&user).Error
	return &user, ErrorToRepoError(err)
}

func (repo *GormUserRepository) GetByFirstnameAndLastname(firstname string, lastname string) ([]domain.User, error) {
	var users []domain.User
	err := repo.Db.
		Where("firstname LIKE ?", "%"+firstname+"%").
		Or("lastname LIKE ?", "%"+lastname+"%").
		Find(&users).Error
	return users, ErrorToRepoError(err)
}

func (repo *GormUserRepository) Create(user *domain.User) error {
	err := repo.Db.Create(user).Error
	return ErrorToRepoError(err)
}

func (repo *GormUserRepository) Update(user *domain.User, updates map[string]interface{}) error {
	err := repo.Db.Model(user).Updates(updates).Error
	return ErrorToRepoError(err)
}

func (repo *GormUserRepository) Delete(id int) error {
	err := repo.Db.Delete(&domain.User{}, id).Error
	return ErrorToRepoError(err)
}
