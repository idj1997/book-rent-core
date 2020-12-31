package domain

type RepoErrorType int

const (
	Unknown              RepoErrorType = 0
	NotFound             RepoErrorType = 1
	InvalidField         RepoErrorType = 2
	UniqueConstraint     RepoErrorType = 3
	ForeignKeyConstraint RepoErrorType = 4
)

type RepoError struct {
	Message string
	Type    RepoErrorType
}

func (e RepoError) Error() string {
	return e.Message
}
