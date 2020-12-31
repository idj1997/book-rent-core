package service

import "github.com/idj1997/book-rent-core/domain"

var NilServiceErrPtr *ServiceError

type ServiceErrorType int

const (
	Unknown               ServiceErrorType = 0
	NotFound              ServiceErrorType = 1
	AlreadyExist          ServiceErrorType = 2
	InvalidArguments      ServiceErrorType = 3
	NotEnoughBooksOnStock ServiceErrorType = 4
	BookAlreadyReturned   ServiceErrorType = 5
	ActiveBookRents       ServiceErrorType = 6
)

type ServiceError struct {
	Type    ServiceErrorType
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func RepoErrorToServiceError(err error) error {
	repoErr := err.(*domain.RepoError)
	var errType ServiceErrorType

	if repoErr != nil {
		if repoErr.Type == domain.UniqueConstraint {
			errType = AlreadyExist
		} else if repoErr.Type == domain.NotFound {
			errType = NotFound
		} else {
			errType = Unknown
		}
		return &ServiceError{Type: errType}
	}
	return nil
}
