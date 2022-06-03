package data

import "github.com/hrand1005/training-notebook/models"

// MockUserDB is my crack at manually implementing a Mock interface for testing
type MockUserDB struct {
	AddUserStub    func(*models.User) (models.UserID, error)
	UsersStub      func() ([]*models.User, error)
	UserByIDStub   func(models.UserID) (*models.User, error)
	UpdateUserStub func(models.UserID, *models.User) error
	DeleteUserStub func(models.UserID) error
	CloseStub      func() error
}

func (m *MockUserDB) AddUser(s *models.User) (models.UserID, error) {
	return m.AddUserStub(s)
}

func (m *MockUserDB) Users() ([]*models.User, error) {
	return m.UsersStub()
}

func (m *MockUserDB) UserByID(id models.UserID) (*models.User, error) {
	return m.UserByIDStub(id)
}

func (m *MockUserDB) UpdateUser(id models.UserID, u *models.User) error {
	return m.UpdateUserStub(id, u)
}

func (m *MockUserDB) DeleteUser(id models.UserID) error {
	return m.DeleteUserStub(id)
}

func (m *MockUserDB) Close() error {
	return m.CloseStub()
}
