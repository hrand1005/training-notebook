package mocks

import (
	"github.com/hrand1005/training-notebook/internal/app"
)

type UserService struct {
	CreateStub func(u *app.User) (app.UserID, error)
}

func (s *UserService) Create(u *app.User) (app.UserID, error) {
	return s.CreateStub(u)
}
