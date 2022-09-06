package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoHandle struct {
	db     *mongo.Database
	client *mongo.Client
	ctx    context.Context
}

// New is the main constructor for the mongo package.
// It returns a handle which can be used to create data stores using the other
// functions offered by this package. Initializes database connection with
// a fresh context. Creates a database on the server with the db name.
// Returns error if the provided uri is invalid.
func New(uri, db string) (*mongoHandle, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db with uri %s, err: %v", uri, err)
	}

	return &mongoHandle{
		db:     client.Database(db),
		client: client,
		ctx:    ctx,
	}, nil
}

// Close cleans up all outstanding database resources and disconnects the client.
func (h *mongoHandle) Close() {
	h.client.Disconnect(h.ctx)
}

// Delete drops the database for the mongoHandle.
func (h *mongoHandle) Delete() {
	h.db.Drop(h.ctx)
}
