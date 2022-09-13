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
	UpdateByID(UserID, *User) error
	FindByID(UserID) (*User, error)
	DeleteByID(UserID) error
}

type UserService interface {
	Create(*User) (UserID, error)
	ReadByID(UserID) (*User, error)
	UpdateByID(UserID, *User) error
	DeleteByID(UserID) error
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
	if err := ValidateEntity(u); err != nil {
		return InvalidUserID, fmt.Errorf("UserService.Create.ValidateEntity: %w", err)
	}
	return s.store.Insert(u)
}

func (s *userService) ReadByID(id UserID) (*User, error) {
	return s.store.FindByID(id)
}

func (s *userService) UpdateByID(id UserID, u *User) error {
	if err := ValidateEntity(u); err != nil {
		return fmt.Errorf("UserService.UpdateByID.ValidateEntity: %w", err)
	}
	return s.store.UpdateByID(id, u)
}

func (s *userService) DeleteByID(id UserID) error {
	return s.store.DeleteByID(id)
}
