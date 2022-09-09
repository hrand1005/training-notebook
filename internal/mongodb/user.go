package mongodb

import (
	"context"
	"errors"
	"fmt"

	"github.com/hrand1005/training-notebook/internal/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string             `bson:"first-name,omitempty"`
	LastName     string             `bson:"last-name,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PasswordHash string             `bson:"password-hash,omitempty"`
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

// Insert saves the user model in the user collection.
// If successful, returns UserID, else returns error.
func (s *userStore) Insert(u *app.User) (app.UserID, error) {
	doc := userToDocument(u)
	res, err := s.coll.InsertOne(s.ctx, doc)
	if err != nil {
		return app.InvalidUserID, fmt.Errorf("failed to insert user: %v", err)
	}

	return userIDFromResult(res), nil
}

// FindByID retrieves the user with the provided id from the user collection.
// If successful, returns a user, otherwise returns nil and error.
func (s *userStore) FindByID(id app.UserID) (*app.User, error) {
	docID, _ := primitive.ObjectIDFromHex(string(id))
	res := s.coll.FindOne(s.ctx, bson.M{"_id": docID})

	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, app.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", app.ErrServiceFailure, err)
	}

	var doc UserDocument
	if err := res.Decode(&doc); err != nil {
		return nil, fmt.Errorf("%w: %v", app.ErrServiceFailure, err)
	}

	return documentToUser(&doc), nil
}

// UpdateByID updates the user identified by the provided id with the fields in the
// provided user struct. If successful, returns nil, otherwise returns error.
func (s *userStore) UpdateByID(id app.UserID, u *app.User) error {
	docID, err := primitive.ObjectIDFromHex(string(id))
	doc := userToDocument(u)
	// TODO: Is it possible to update only the provided fields?
	res, err := s.coll.UpdateByID(s.ctx, docID, bson.M{
		"$set": doc,
	})

	if err != nil {
		return fmt.Errorf("%w: %v", app.ErrServiceFailure, err)
	}

	if res.MatchedCount == 0 {
		return app.ErrNotFound
	}

	return nil
}

/* --- INTERNAL USE ONLY --- */

func documentToUser(d *UserDocument) *app.User {
	return &app.User{
		ID:           app.UserID(d.ID.Hex()),
		FirstName:    d.FirstName,
		LastName:     d.LastName,
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
	}
}

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
