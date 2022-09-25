package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/hrand1005/training-notebook/internal/app"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

type setStore struct {
	coll *mongo.Collection
	ctx  context.Context
}

const SetCollection = "sets"

func NewSetStore(h *mongoHandle) app.SetStore {
	return &setStore{
		coll: h.db.Collection(SetCollection),
		ctx:  h.ctx,
	}
}

type SetDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	OwnerID   app.UserID         `bson:"owner-id,omitempty"`
	Movement  string             `bson:"movement,omitempty"`
	Intensity float64            `bson:"intensity,omitempty"`
	Volume    int                `bson:"volume,omitempty"`
}

func (s *setStore) Insert(set *app.Set) (app.SetID, error) {
	doc := setToDocument(set)
	res, err := s.coll.InsertOne(s.ctx, doc)
	if err != nil {
		return app.InvalidSetID, fmt.Errorf("failed to insert set: %v", err)
	}

	return setIDFromResult(res), nil
}

func (s *setStore) UpdateByID(id app.SetID, set *app.Set) error {
	log.Println("NOT IMPLEMENTED")
	return nil
}

func (s *setStore) FindByID(id app.SetID) (*app.Set, error) {
	log.Println("NOT IMPLEMENTED")
	return nil, nil
}

func (s *setStore) DeleteByID(id app.SetID) error {
	log.Println("NOT IMPLEMENTED")
	return nil
}

/* --- INTERNAL USE ONLY --- */

func documentToSet(s *SetDocument) *app.Set {
	return &app.Set{
		ID:        app.SetID(s.ID.Hex()),
		OwnerID:   s.OwnerID,
		Movement:  s.Movement,
		Intensity: s.Intensity,
		Volume:    s.Volume,
	}
}

func setToDocument(s *app.Set) *SetDocument {
	return &SetDocument{
		OwnerID:   s.OwnerID,
		Movement:  s.Movement,
		Intensity: s.Intensity,
		Volume:    s.Volume,
	}
}

func setIDFromResult(r *mongo.InsertOneResult) app.SetID {
	return app.SetID(r.InsertedID.(primitive.ObjectID).Hex())
}
