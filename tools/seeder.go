package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

var (
	conf = flag.String("config", "", "config file defining the database to be populated")
)

const (
	NumUsers = 50
	NumSets  = 1000
)

func main() {

	log.SetFlags(log.Ltime)
	log.SetPrefix("SEEDER: ")

	flag.Parse()
	if *conf == "" {
		log.Fatal("--config must be set")
	}

	dbConf, err := loadDBConfig(*conf)
	if err != nil {
		log.Fatalf("loadDBConfig: %v", err)
	}

	// connect to database using config values
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConf.URI))
	if err != nil {
		log.Fatalf("failed to create mongo db client: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("failed to ping mongo db: %v", err)
	}

	db := client.Database(dbConf.Name)

	// prepare for random data generation and seed the database
	rand.Seed(time.Now().UnixNano())

	log.Printf("Populating collection %q with %v users...", UserCollection, NumUsers)
	usersResult, err := seedUsers(ctx, db)
	if err != nil {
		log.Fatalf("seedUsers: %v", err)
	}

	log.Printf("Populating collection %q with %v sets...", SetCollection, NumSets)
	_, err = seedSets(ctx, db, usersResult.InsertedIDs)
	if err != nil {
		log.Fatalf("seedSets: %v", err)
	}

	log.Printf("Seeding of %q is complete, exiting...", dbConf.Name)
}

// TODO: Move to shared package when required by the server
type DBConfig struct {
	Name string `yaml:"database-name"`
	URI  string `yaml:"mongodb-uri"`
}

// loadDBConfig loads relevant mongodb information from the provided
// config file path.
func loadDBConfig(file string) (DBConfig, error) {
	f, err := os.Open(file)
	if err != nil {
		return DBConfig{}, fmt.Errorf("opening file: %s, err: %v", file, err)
	}
	defer f.Close()

	var dbConf DBConfig
	d := yaml.NewDecoder(f)
	if err := d.Decode(&dbConf); err != nil {
		return DBConfig{}, fmt.Errorf("decoding config: %v", err)
	}

	return dbConf, nil
}

// TODO: Move models to shared package
const UserCollection = "users"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string             `bson:"first-name"`
	LastName     string             `bson:"last-name"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password-hash"`
}

// seedUsers populates the provided database with NumUsers semi-random users.
func seedUsers(ctx context.Context, db *mongo.Database) (*mongo.InsertManyResult, error) {
	collection := db.Collection(UserCollection)
	collection.DeleteMany(ctx, bson.D{})

	var users []interface{}
	for i := 0; i < NumUsers; i++ {
		newUser := User{
			FirstName:    gofakeit.FirstName(),
			LastName:     gofakeit.LastName(),
			Email:        gofakeit.Email(),
			PasswordHash: generateHash(gofakeit.Password(true, true, true, true, true, 32)),
		}
		users = append(users, newUser)
	}

	return collection.InsertMany(ctx, users)
}

const SetCollection = "sets"

type Set struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user-id"`
	Movement  string             `bson:"movement"`
	Volume    int                `bson:"volume"`
	Intensity float64            `bson:"intensity"`
	Date      primitive.DateTime `bson:"date"`
}

// seedSets populates the provided database with NumSets semi-random sets.
// Each set belongs to one of the provided userIDs.
func seedSets(ctx context.Context, db *mongo.Database, userIDs []interface{}) (*mongo.InsertManyResult, error) {
	collection := db.Collection(SetCollection)
	collection.DeleteMany(ctx, bson.D{})

	var sets []interface{}
	for i := 0; i < NumSets; i++ {
		// assign set a random existing user
		userID := userIDs[rand.Intn(len(userIDs))]
		newSet := Set{
			UserID:    userID.(primitive.ObjectID),
			Movement:  gofakeit.VerbAction(),
			Volume:    rand.Intn(100),
			Intensity: rand.Float64() * 100,
			Date:      primitive.NewDateTimeFromTime(gofakeit.Date()),
		}
		sets = append(sets, newSet)
	}

	return collection.InsertMany(ctx, sets)
}

func generateHash(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashBytes)
}
