package models

type Credentails struct {
	UserID   UserID `json:"user-id"`
	Password string `json:"password,omitempty"`
}
