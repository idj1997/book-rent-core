package domain

import "gorm.io/gorm"

type UserType int

const (
	ADMIN    UserType = 0
	CUSTOMER UserType = 1
)

type User struct {
	gorm.Model
	Firstname string
	Lastname  string
	Email     string `gorm:"unique"`
	Password  string `gorm:"not null"`
	Type      UserType
}

type UserRepository interface {
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByFirstnameAndLastname(firstname string, lastname string) ([]User, error)
	Create(user *User) error
	Update(user *User, updates map[string]interface{}) error
	Delete(id int) error
}

type UserService interface {
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	GetByFirstnameAndLastname(firstname string, lastname string) ([]User, error)
	Create(user *User) error
	Delete(id int) error
}
