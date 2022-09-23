package mocks

import (
	"github.com/hrand1005/training-notebook/internal/app"
)

type SetService struct {
	CreateStub     func(*app.Set) (app.SetID, error)
	ReadByIDStub   func(app.SetID) (*app.Set, error)
	UpdateByIDStub func(app.SetID, *app.Set) error
	DeleteByIDStub func(app.SetID) error
}

func (s *SetService) Create(u *app.Set) (app.SetID, error) {
	return s.CreateStub(u)
}

func (s *SetService) ReadByID(id app.SetID) (*app.Set, error) {
	return s.ReadByIDStub(id)
}

func (s *SetService) UpdateByID(id app.SetID, u *app.Set) error {
	return s.UpdateByIDStub(id, u)
}

func (s *SetService) DeleteByID(id app.SetID) error {
	return s.DeleteByIDStub(id)
}
