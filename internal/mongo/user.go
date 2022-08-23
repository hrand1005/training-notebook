package mongo

import (
	"context"
	"fmt"

	"github.com/hrand1005/training-notebook/internal/app"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDocument struct {
	FirstName    string `bson:"first-name"`
	LastName     string `bson:"last-name"`
	Email        string `bson:"email"`
	PasswordHash string `bson:"password-hash"`
}

type userStore struct {
	coll *mongo.Collection
	ctx  context.Context
}

const UserCollection = "users"

// NewUserStore creates a user collection in the database in the
// provided mongo handle. Returns a UserStore for data operations.
func NewUserStore(h *mongoHandle) app.UserStore {
	return &userStore{
		coll: h.db.Collection(UserCollection),
		ctx:  h.ctx,
	}
}

// Insert saves the user model in the UserStore.
// If successful, returns UserID, else returns error.
func (s *userStore) Insert(u *app.User) (app.UserID, error) {
	doc := userToDocument(u)
	res, err := s.coll.InsertOne(s.ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert user: %v", err)
	}

	return userIDFromResult(res), nil
}

/* --- INTERNAL USE ONLY --- */

func userToDocument(u *app.User) *UserDocument {
	return &UserDocument{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}
}

func userIDFromResult(r *mongo.InsertOneResult) app.UserID {
	return app.UserID(r.InsertedID.(primitive.ObjectID).Hex())
}
