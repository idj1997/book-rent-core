package service

import (
	"github.com/idj1997/book-rent-core/domain"
	"time"
)

type RentDetailsService struct {
	RentRepo domain.RentDetailsRepository
	BookRepo domain.BookRepository
}

func (r *RentDetailsService) GetByID(id int) (*domain.RentDetails, error) {
	rent, err := r.RentRepo.GetByID(id)
	return rent, RepoErrorToServiceError(err)
}

func (r *RentDetailsService) RentBook(rent *domain.RentDetails) error {
	book, getBookErr := r.BookRepo.GetByID(rent.BookID)
	if getBookErr != domain.NilRepoErrPtr {
		return RepoErrorToServiceError(getBookErr)
	}

	if book.Stock <= 0 {
		return &ServiceError{Type: NotEnoughBooksOnStock}
	}

	rent.CreatedAt = time.Now()
	rent.ReturnDeadline = time.Now().Add(30 * 24 * time.Hour)
	createRentErr := r.RentRepo.Create(rent)
	if createRentErr != domain.NilRepoErrPtr {
		return RepoErrorToServiceError(createRentErr)
	}

	stockUpdate := make(map[string]interface{})
	stockUpdate["stock"] = book.Stock - 1

	updateStockErr := r.BookRepo.Update(book, stockUpdate)
	if updateStockErr != domain.NilRepoErrPtr {
		return RepoErrorToServiceError(updateStockErr)
	}

	return NilServiceErrPtr
}

func (r *RentDetailsService) ReturnBook(rentDetailsID int) error {
	rent, getRentErr := r.RentRepo.GetByID(rentDetailsID)
	if getRentErr != domain.NilRepoErrPtr {
		return RepoErrorToServiceError(getRentErr)
	}

	if rent.Status == domain.RETURNED {
		return &ServiceError{Type: BookAlreadyReturned}
	}

	bookUpdates := make(map[string]interface{})
	bookUpdates["stock"] = rent.Book.Stock + 1

	rentUpdates := make(map[string]interface{})
	rentUpdates["status"] = domain.RETURNED

	updateRentErr := r.RentRepo.Update(rent, rentUpdates)
	if updateRentErr != domain.NilRepoErrPtr {
		return RepoErrorToServiceError(updateRentErr)
	}

	updateBookErr := r.BookRepo.Update(&rent.Book, bookUpdates)
	if updateBookErr != domain.NilRepoErrPtr {
		return RepoErrorToServiceError(updateBookErr)
	}

	return NilServiceErrPtr
}

func (r *RentDetailsService) GetByUser(userID int) ([]domain.RentDetails, error) {
	rents, err := r.RentRepo.GetByUser(userID)
	return rents, RepoErrorToServiceError(err)
}

func (r *RentDetailsService) GetByBook(bookID int) ([]domain.RentDetails, error) {
	rents, err := r.RentRepo.GetByBook(bookID)
	return rents, RepoErrorToServiceError(err)
}

func (r *RentDetailsService) GetByStatus(status domain.RentDetailsStatus) ([]domain.RentDetails, error) {
	rents, err := r.RentRepo.GetByStatus(status)
	return rents, RepoErrorToServiceError(err)
}

func (r *RentDetailsService) UpdateToExpired() error {
	stream := make(chan domain.RentDetails)
	updates := make(map[string]interface{})
	updates["status"] = domain.EXPIRED

	go r.RentRepo.RentDetailsIterator(stream)
	for rent := range stream {
		if rent.Status == domain.RENTED && rent.ReturnDeadline.Before(time.Now()) {
			err := r.RentRepo.Update(&rent, updates)
			if err != nil {
				return RepoErrorToServiceError(err)
			}
		}
	}

	return nil
}
