package app

type UserID string

type User struct {
	ID           UserID
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
}

type UserStore interface {
	Insert(*User) (UserID, error)
	FindByID(UserID) (*User, error)
}

type UserService interface {
	Create(*User) (UserID, error)
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
	return s.store.Insert(u)
}

func (s *userService) ReadByID(id UserID) (*User, error) {
	return s.store.FindByID(id)
}
