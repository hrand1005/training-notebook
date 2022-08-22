package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: move to models?
type User struct {
	// TODO: decide ID field type
	ID           string `bson:"_id,omitempty"`
	FirstName    string `bson:"first-name"`
	LastName     string `bson:"last-name"`
	Email        string `bson:"email"`
	PasswordHash string `bson:"password-hash"`
}

type UserStore interface {
	Insert(u *User) (string, error)
}

type userStore struct {
	coll *mongo.Collection
	ctx  context.Context
}

const UserCollection = "users"

func NewUserStore(h *mongoHandle) UserStore {
	return &userStore{
		coll: h.db.Collection(UserCollection),
		ctx:  h.ctx,
	}
}

func (s userStore) Insert(u *User) (string, error) {
	res, err := s.coll.InsertOne(s.ctx, u)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %v", err)
	}

	objectID := res.InsertedID.(primitive.ObjectID)

	return objectID.String(), nil
}
