package data

import "github.com/hrand1005/training-notebook/models"

// MockUserDB is my crack at manually implementing a Mock interface for testing
type MockUserDB struct {
	AddUserStub func(*models.User) (models.UserID, error)
	UsersStub   func() ([]*models.User, error)
}

func (m *MockUserDB) AddUser(s *models.User) (models.UserID, error) {
	return m.AddUserStub(s)
}

func (m *MockUserDB) Users() ([]*models.User, error) {
	return m.UsersStub()
}
