package app

import (
	"fmt"
)

type SetID string

const InvalidSetID = "0"

type Set struct {
	ID        SetID
	OwnerID   UserID
	Movement  string  `validate:"required,min=2,max=32"`
	Intensity float64 `validate:"required,min=0,max=100"`
	Volume    int     `validate:"required,min=1,max=100"`
}

type SetStore interface {
	Insert(*Set) (SetID, error)
	UpdateByID(SetID, *Set) error
	FindByID(SetID) (*Set, error)
	DeleteByID(SetID) error
}

type SetService interface {
	Create(*Set) (SetID, error)
	ReadByID(SetID) (*Set, error)
	UpdateByID(SetID, *Set) error
	DeleteByID(SetID) error
}

type setService struct {
	store SetStore
}

func NewSetService(store SetStore) SetService {
	return &setService{
		store: store,
	}
}

func (s *setService) Create(u *Set) (SetID, error) {
	if err := ValidateEntity(u); err != nil {
		return InvalidSetID, fmt.Errorf("SetService.Create.ValidateEntity: %w", err)
	}
	return s.store.Insert(u)
}

func (s *setService) ReadByID(id SetID) (*Set, error) {
	return s.store.FindByID(id)
}

func (s *setService) UpdateByID(id SetID, u *Set) error {
	if err := ValidateEntity(u); err != nil {
		return fmt.Errorf("SetService.UpdateByID.ValidateEntity: %w", err)
	}
	return s.store.UpdateByID(id, u)
}

func (s *setService) DeleteByID(id SetID) error {
	return s.store.DeleteByID(id)
}
