package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping mongo db: %v", err)
	}

	mode := os.Getenv("MODE")
	if mode == "test" {
		// use client to populate the database with random data
		if err := SeedDB(ctx, client, db); err != nil {
			log.Fatalf("Failed to seed database: %v", err)
		}
	}

	app := fiber.New()

	app.Get("/bson-map", func(c *fiber.Ctx) error {
		var results []bson.M
		cur, err := client.Database(db).Collection("test").Find(ctx, bson.D{})
		if err != nil {
			return fmt.Errorf("get request on '/bson-map': %v", err)
		}
		if err = cur.All(ctx, &results); err != nil {
			return fmt.Errorf("get request on '/bson-map': %v", err)
		}
		fmt.Printf("Results as bson mapping:\n%+v", results)

		return c.JSON(results)
	})

	app.Get("/struct", func(c *fiber.Ctx) error {
		type result struct {
			ID  primitive.ObjectID `bson:"_id"`
			Key string             `bson:"key,omitempty"`
		}
		var results []result
		cur, err := client.Database(db).Collection("test").Find(ctx, bson.D{})
		if err != nil {
			return fmt.Errorf("get request on '/struct': %v", err)
		}
		defer cur.Close(ctx)
		for cur.Next(ctx) {
			r := result{}
			if err := cur.Decode(&r); err != nil {
				return fmt.Errorf("get request on '/struct': %v", err)
			}
			results = append(results, r)
		}
		fmt.Printf("Results in struct slice:\n%+v", results)

		return c.JSON(results)
	})

	app.Listen(":5000")
}

// TODO: Make semi-random and useful for testing
func SeedDB(ctx context.Context, client *mongo.Client, db string) error {
	collection := client.Database(db).Collection("test")
	_, err := collection.InsertOne(ctx, bson.M{"key": "value"})
	return err
}
