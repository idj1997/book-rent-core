package repository

import (
	"errors"
	"github.com/idj1997/book-rent-core/domain"
	"strings"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ErrorToRepoError(err error) *domain.RepoError {
	if err != nil {
		var errType domain.RepoErrorType
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errType = domain.NotFound
		} else if strings.Contains(err.Error(), "23505") {
			errType = domain.UniqueConstraint
		} else if strings.Contains(err.Error(), "42703") {
			errType = domain.InvalidField
		} else if strings.Contains(err.Error(), "23503") {
			errType = domain.ForeignKeyConstraint
		} else {
			log.Errorf("unknown/unexpected error: %v", err)
			errType = domain.Unknown
		}
		return &domain.RepoError{Type: errType, Message: err.Error()}
	}
	return nil
}
