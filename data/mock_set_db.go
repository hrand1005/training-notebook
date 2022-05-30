package data

// MockSetDB is my crack at manually implementing a Mock interface for testing
type MockSetDB struct {
	AddSetStub: func(s *Set) (SetID, error)
	SetsStub: func() ([]*Set, error)
	SetByIDStub: func(id SetID) (*Set, error)
	UpdateSetStub: func(id SetID, s *Set) error
	DeleteSetStub: func(id SetID) error
	CloseStub: func() error
}

func (m *MockSetDB) AddSet(s *Set) (SetID, error) {
	return m.AddSetStub(s)
}

func (m *MockSetDB) Sets() ([]*Set, error) {
	return m.SetsStub()
}

func (m *MockSetDB) SetByID(id SetID) (*Set, error) {
	return m.SetByIDStub(id)
}

func (m *MockSetDB) UpdateSet(id SetID, s *Set) error {
	return m.UpdateSet(id, s)
}

func (m *MockSetDB) DeleteSet(id SetID) error {
	return m.DeleteSet(id)
}

func (m *MockSetDB) Close() error {
	return m.Close()
}