package app

import (
	"fmt"
)

type UserID string

const InvalidUserID = "0"

type User struct {
	ID           UserID
	FirstName    string `validate:"required,min=2,max=32"`
	LastName     string `validate:"required,min=2,max=32"`
	Email        string `validate:"required,email,min=6,max=32"`
	PasswordHash string
}

type UserStore interface {
	Insert(*User) (UserID, error)
	FindByID(UserID) (*User, error)
}

type UserService interface {
	Create(*User) (UserID, error)
	ReadByID(id UserID) (*User, error)
}

type userService struct {
	store UserStore
}

func NewUserService(store UserStore) UserService {
	return &userService{
		store: store,
	}
}

func (s *userService) Create(u *User) (UserID, error) {
	if errors := ValidateEntity(u); errors != nil {
		return InvalidUserID, fmt.Errorf("UserService.Create.ValidateEntity: %w", errors)
	}
	return s.store.Insert(u)
}

func (s *userService) ReadByID(id UserID) (*User, error) {
	return s.store.FindByID(id)
}
