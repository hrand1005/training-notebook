package app

type SetID string

type Set struct {
	ID        SetID
	OwnerID   UserID
	Movement  string
	Intensity float64
	Volume    int
}
