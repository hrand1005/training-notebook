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
}

type UserService interface {
	Create(UserID) (*User, error)
}
