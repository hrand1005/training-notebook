package data

// MockUserDB is my crack at manually implementing a Mock interface for testing
type MockUserDB struct {
	AddUserStub func(*User) (UserID, error)
	UsersStub   func() ([]*User, error)
}

func (m *MockUserDB) AddUser(s *User) (UserID, error) {
	return m.AddUserStub(s)
}

func (m *MockUserDB) Users() ([]*User, error) {
	return m.UsersStub()
}
