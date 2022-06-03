package models

// UserID is the unique identifier for a user of the application
type UserID int

// User defines the model of the user resource
type User struct {
	ID   UserID `json:"id"`
	Name string `json:"name"`
}
