package mocks

import (
	"github.com/hrand1005/training-notebook/internal/app"
)

type UserService struct {
	CreateStub     func(*app.User) (app.UserID, error)
	ReadByIDStub   func(app.UserID) (*app.User, error)
	UpdateByIDStub func(app.UserID, *app.User) error
}

func (s *UserService) Create(u *app.User) (app.UserID, error) {
	return s.CreateStub(u)
}

func (s *UserService) ReadByID(id app.UserID) (*app.User, error) {
	return s.ReadByIDStub(id)
}

func (s *UserService) UpdateByID(id app.UserID, u *app.User) error {
	return s.UpdateByIDStub(id, u)
}
