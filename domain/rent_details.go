package domain

import (
	"time"

	"gorm.io/gorm"
)

type RentDetailsStatus int

const (
	RENTED   RentDetailsStatus = 0
	RETURNED RentDetailsStatus = 1
	EXPIRED  RentDetailsStatus = 2
)

type RentDetails struct {
	gorm.Model
	UserID         int               `gorm:"not null"`
	BookID         int               `gorm:"not null"`
	Status         RentDetailsStatus `gorm:"default:0"`
	ReturnedAt     time.Time
	ReturnDeadline time.Time
	User           User
	Book           Book
}

type RentDetailsRepository interface {
	GetByID(id int) (*RentDetails, error)
	Create(rent *RentDetails) error
	Update(rent *RentDetails, updates map[string]interface{}) error
	UpdateAssociations(rent *RentDetails, updates map[string]interface{}) error
	GetByUser(userID int) ([]RentDetails, error)
	GetByBook(bookID int) ([]RentDetails, error)
	GetByStatus(status RentDetailsStatus) ([]RentDetails, error)
	RentDetailsIterator(chan RentDetails)
}

type RentDetailsService interface {
	GetByID(id int) (*RentDetails, error)
	RentBook(rent *RentDetails) error
	ReturnBook(rentDetailsID int) error
	GetByUser(userID int) ([]RentDetails, error)
	GetByBook(bookID int) ([]RentDetails, error)
	GetByStatus(status RentDetailsStatus) ([]RentDetails, error)
	UpdateToExpired() error
}
