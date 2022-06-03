package models

// UserID is the unique identifier for a user of the application
type UserID int

// User defines the model of the user resource
type User struct {
	ID   UserID `json:"id"`
	Name string `json:"name"`
}

// UsersEqual returns true if all non-id fields of the user are equal, and false otherwise.
func UsersEqual(user1, user2 *User) bool {
	if user1.Name != user2.Name {
		return false
	}
	return true
}

// ContainsUser checks if the slice of users contains the provided user.
func ContainsUser(users []*User, s *User) bool {
	for _, v := range users {
		if UsersEqual(s, v) {
			return true
		}
	}
	return false
}
