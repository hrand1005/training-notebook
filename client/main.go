package main

import (
	"context"
	"log"

	"github.com/hrand1005/training-notebook/set"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(":8080", opts...)
	if err != nil {
		log.Fatalf("Failed to get conn with dial: %#v", err)
	}
	defer conn.Close()
	cl := set.NewSetClient(conn)
	ctx := context.Background()

	newSet := &set.CreateRequest{
		Id: 1,
		Name: "Curl",
		Reps: 20,
		Intensity: 65.0,
	}

	log.Printf("Creating set:\n%#v", newSet)
	resp, err := cl.Create(ctx, newSet)
	if err != nil {
		log.Fatalf("Encountered error creating new set: %v", err)
	}
	log.Printf("Got create response:\n%#v", resp)
	return
}