package data

import "github.com/hrand1005/training-notebook/models"

// MockSetDB is my crack at manually implementing a Mock interface for testing
type MockSetDB struct {
	AddSetStub           func(s *models.Set) (models.SetID, error)
	SetsStub             func() ([]*models.Set, error)
	SetsByUserIDStub     func(models.UserID) ([]*models.Set, error)
	SetByIDStub          func(id models.SetID) (*models.Set, error)
	SetByIDForUserStub   func(models.SetID, models.UserID) (*models.Set, error)
	UpdateSetStub        func(id models.SetID, s *models.Set) error
	UpdateSetForUserStub func(setID models.SetID, userID models.UserID, s *models.Set) error
	DeleteSetStub        func(id models.SetID) error
	DeleteSetForUserStub func(setID models.SetID, userID models.UserID) error
	CloseStub            func() error
}

func (m *MockSetDB) AddSet(s *models.Set) (models.SetID, error) {
	return m.AddSetStub(s)
}

func (m *MockSetDB) Sets() ([]*models.Set, error) {
	return m.SetsStub()
}

func (m *MockSetDB) SetsByUserID(id models.UserID) ([]*models.Set, error) {
	return m.SetsByUserIDStub(id)
}

func (m *MockSetDB) SetByID(id models.SetID) (*models.Set, error) {
	return m.SetByIDStub(id)
}

func (m *MockSetDB) SetByIDForUser(setID models.SetID, userID models.UserID) (*models.Set, error) {
	return m.SetByIDForUserStub(setID, userID)
}

func (m *MockSetDB) UpdateSet(id models.SetID, s *models.Set) error {
	return m.UpdateSetStub(id, s)
}

func (m *MockSetDB) UpdateSetForUser(setID models.SetID, userID models.UserID, s *models.Set) error {
	return m.UpdateSetForUserStub(setID, userID, s)
}

func (m *MockSetDB) DeleteSet(id models.SetID) error {
	return m.DeleteSetStub(id)
}

func (m *MockSetDB) DeleteSetForUser(setID models.SetID, userID models.UserID) error {
	return m.DeleteSetForUserStub(setID, userID)
}

func (m *MockSetDB) Close() error {
	return m.CloseStub()
}
