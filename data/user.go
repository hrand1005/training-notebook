package data

type User struct {
	ID     int    `json:"id"`
	UserID string `json:"userID"`
	Name   string `json:"name"`
}

var users = []*User{
	{
		ID:     1,
		UserID: "hrand",
		Name:   "Herbie",
	},
}

/*func UserByID(userID string) (*User, error) {
	for _, u := range users {
		if u.UserID == userID {
			return u, nil
		}
	}

	return nil, ErrNotFound
}*/

func UserByUserID(userID string) (*User, error) {
	for _, u := range users {
		if u.UserID == userID {
			return u, nil
		}
	}

	return nil, ErrNotFound
}
