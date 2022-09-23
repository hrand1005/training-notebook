package mongodb

import (
	"context"

	"github.com/hrand1005/training-notebook/internal/app"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
)

type SetStore interface {
}

type setStore struct {
	coll *mongo.Collection
	ctx  context.Context
}

const SetCollection = "sets"

func NewSetStore(h *mongoHandle) SetStore {
	return &setStore{
		coll: h.db.Collection(SetCollection),
		ctx:  h.ctx,
	}
}

type SetDocument struct {
	ID        app.SetID  `bson:"_id,omitempty"`
	OwnerID   app.UserID `bson:"owner-id,omitempty"`
	Movement  string     `bson:"movement,omitempty"`
	Intensity float64    `bson:"intensity,omitempty"`
	Volume    int        `bson:"volume,omitempty"`
}
