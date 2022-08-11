package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file variables")
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		// TODO: set default / testing values?
		log.Fatalf("MONGODB_URI not set")
	}

	db := os.Getenv("DATABASE")
	if db == "" {
		log.Fatalf("DATABASE not set")
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create mongo db client: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping mongo db: %v", err)
	}

	collection := client.Database(db).Collection("test")
	res, err := collection.InsertOne(ctx, bson.M{"key": "value"})
	if err != nil {
		log.Fatalf("Failed to insert test collection: %v", err)
	}

	id := res.InsertedID
	log.Printf("Inserted collection with ID: %v", id)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		var result bson.M
		err := collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&result)

		if err != nil {
			return fmt.Errorf("get request on '/': %v", err)
		}
		return c.JSON(result)
	})

	app.Listen(":5000")
}
