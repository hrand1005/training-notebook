package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoHandle struct {
	client *mongo.Client
	ctx    context.Context
}

// New is the main constructor for the mongo package.
// It returns a handle which can be used to create data stores using the other
// functions offered by this package. Initializes database connection with
// a fresh context.
// Returns error if the provided uri is invalid.
func New(uri string) (*mongoHandle, error) {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db with uri %s, err: %v", uri, err)
	}

	return &mongoHandle{
		client: client,
		ctx:    ctx,
	}, nil
}

// Close cleans up all outstanding database resources and disconnects the client.
func (h *mongoHandle) Close() {
	h.client.Disconnect(h.ctx)
}
