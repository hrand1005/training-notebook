package data

import (
	"errors"
)

var ErrInvalidUserID = errors.New("user ID is not valid")

var ErrUserIDExists = errors.New("user ID is already in user")

type User struct {
	UserID string `json:"userID"`
	Name   string `json:"name"`
}

var users = []*User{
	{
		UserID: "hrand",
		Name:   "Herbie",
	},
}

// Replace with DB logic
func AddUser(u *User) error {
	if err := checkUserID(u.UserID); err != nil {
		return err
	}
	users = append(users, u)
	return nil
}

func UserByUserID(userID string) (*User, error) {
	for _, u := range users {
		if u.UserID == userID {
			return u, nil
		}
	}

	return nil, ErrNotFound
}

func UpdateUser(userID string, u *User) error {
	for i, u := range users {
		if userID == u.UserID {
			u.UserID = userID
			users[i] = u
			return nil
		}
	}

	return ErrNotFound
}

func DeleteUser(userID string) error {
	for i, u := range users {
		if u.UserID == userID {
			users = append(users[:i], users[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}

func checkUserID(userID string) error {
	if userID == "" {
		return ErrInvalidUserID
	}

	for _, u := range users {
		if u.UserID == userID {
			return ErrUserIDExists
		}
	}

	return nil
}
