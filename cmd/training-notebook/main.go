package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/hrand1005/training-notebook/internal/config"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var configPath = flag.String("config", "", "Path to file containing server configs")

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load .env file variables")
	}

	flag.Parse()
	conf, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load server configs from %s, err: %v", *configPath, err)
	}

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.Database.URI))
	if err != nil {
		log.Fatalf("failed to create mongo db client: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping mongo db: %v", err)
	}

	app := fiber.New()

	app.Listen(conf.Server.Port)
}

// FOR REFERENCE
// app.Get("/bson-map", func(c *fiber.Ctx) error {
// 	var results []bson.M
// 	cur, err := client.Database(srvConf.Database.Name).Collection("test").Find(ctx, bson.D{})
// 	if err != nil {
// 		return fmt.Errorf("get request on '/bson-map': %v", err)
// 	}
// 	if err = cur.All(ctx, &results); err != nil {
// 		return fmt.Errorf("get request on '/bson-map': %v", err)
// 	}
// 	fmt.Printf("Results as bson mapping:\n%+v", results)
//
// 	return c.JSON(results)
// })
//
// app.Get("/struct", func(c *fiber.Ctx) error {
// 	type result struct {
// 		ID  primitive.ObjectID `bson:"_id"`
// 		Key string             `bson:"key,omitempty"`
// 	}
// 	var results []result
// 	cur, err := client.Database(srvConf.Database.Name).Collection("test").Find(ctx, bson.D{})
// 	if err != nil {
// 		return fmt.Errorf("get request on '/struct': %v", err)
// 	}
// 	defer cur.Close(ctx)
// 	for cur.Next(ctx) {
// 		r := result{}
// 		if err := cur.Decode(&r); err != nil {
// 			return fmt.Errorf("get request on '/struct': %v", err)
// 		}
// 		results = append(results, r)
// 	}
// 	fmt.Printf("Results in struct slice:\n%+v", results)
//
// 	return c.JSON(results)
// })
