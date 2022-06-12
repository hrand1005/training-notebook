package models

type Credentails struct {
	UID      UserID `json:"user-id"`
	Password string `json:"password,omitempty"`
}
