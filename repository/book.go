package repository

import (
	"github.com/idj1997/book-rent-core/domain"
	"gorm.io/gorm"
)

type GormBookRepository struct {
	Db *gorm.DB
}

func NewGormBookRepository(db *gorm.DB) *GormBookRepository {
	return &GormBookRepository{Db: db}
}

func (repo *GormBookRepository) GetByID(id int) (*domain.Book, error) {
	var book domain.Book
	err := repo.Db.First(&book, id).Error
	return &book, ErrorToRepoError(err)
}

func (repo *GormBookRepository) GetByTitle(title string) ([]domain.Book, error) {
	var books []domain.Book
	err := repo.Db.Where("title LIKE ?", "%"+title+"%").Find(&books).Error
	return books, ErrorToRepoError(err)
}

func (repo *GormBookRepository) Create(book *domain.Book) (uint, error) {
	err := repo.Db.Create(book).Error
	return book.ID, ErrorToRepoError(err)
}

func (repo *GormBookRepository) Update(book *domain.Book, updates map[string]interface{}) error {
	err := repo.Db.Model(book).Updates(updates).Error
	return ErrorToRepoError(err)
}

func (repo *GormBookRepository) Delete(id int) error {
	err := repo.Db.Delete(&domain.Book{}, id).Error
	return ErrorToRepoError(err)
}
