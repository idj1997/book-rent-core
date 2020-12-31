package repository

import (
	"github.com/idj1997/book-rent-core/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GormRentDetailsRepository struct {
	Db *gorm.DB
}

func (g *GormRentDetailsRepository) GetByID(id int) (*domain.RentDetails, error) {
	var rent domain.RentDetails
	err := g.Db.
		Preload(clause.Associations).
		First(&rent, id).Error
	return &rent, ErrorToRepoError(err)
}

func (g *GormRentDetailsRepository) Create(rent *domain.RentDetails) error {
	err := g.Db.Create(rent).Error
	return ErrorToRepoError(err)
}

func (g *GormRentDetailsRepository) Update(rent *domain.RentDetails, updates map[string]interface{}) error {
	err := g.Db.
		Model(rent).               // specify model on which to perform updates
		Omit(clause.Associations). // discard associations updates in rent pointer
		Updates(updates).          // perform updates
		Error
	return ErrorToRepoError(err)
}

// UpdateAssociations function will insert updated associations into rent pointer
func (g *GormRentDetailsRepository) UpdateAssociations(rent *domain.RentDetails, updates map[string]interface{}) error {
	err := g.Db.
		Model(rent).
		Omit(clause.Associations).
		Updates(updates).
		Preload(clause.Associations). // select nested associations in rent
		First(rent, rent.ID).         // fetch
		Error
	return ErrorToRepoError(err)
}

func (g *GormRentDetailsRepository) GetByUser(userID int) ([]domain.RentDetails, error) {
	var rents []domain.RentDetails
	err := g.Db.
		Where("user_id=?", userID).
		Find(&rents).
		Error
	return rents, ErrorToRepoError(err)
}

func (g *GormRentDetailsRepository) GetByBook(bookID int) ([]domain.RentDetails, error) {
	var rents []domain.RentDetails
	err := g.Db.
		Where("book_id=?", bookID).
		Find(&rents).
		Error
	return rents, ErrorToRepoError(err)
}

func (g *GormRentDetailsRepository) GetByStatus(status domain.RentDetailsStatus) ([]domain.RentDetails, error) {
	var rents []domain.RentDetails
	err := g.Db.
		Where("status=?", status).
		Find(&rents).
		Error
	return rents, ErrorToRepoError(err)
}

func (g *GormRentDetailsRepository) RentDetailsIterator(stream chan domain.RentDetails) {
	// close channel
	defer close(stream)

	// create rows struct
	rows, err := g.Db.Model(&domain.RentDetails{}).
		Where("status != ?", domain.RETURNED).
		Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	// iterate and stream results to channel
	for rows.Next() {
		var rent domain.RentDetails
		err := g.Db.ScanRows(rows, &rent)
		if err != nil {

		}
		stream <- rent
	}
}
