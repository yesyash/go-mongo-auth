package gomongoauth

import (
	"context"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitializeDbConnection(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		log.Println("Error connecting to db. Err: ", err.Error())
		return client, err
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Println("Unable to ping to primary. Err: ", err.Error())
		return client, err
	}

	log.Println("Connected to db.")

	return client, nil
}

// Must be used as a defer function to disconnect from db.
func DisconnectDb(client *mongo.Client) error {
	if err := client.Disconnect(context.TODO()); err != nil {
		return err
	}

	log.Println("Disconnecting from db.")
	return nil
}

func validateInput(user User) error {
	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		log.Printf("Error validating user input.")
	}

	return err
}

func getUserByEmail(collection mongo.Collection, ctx context.Context, email string) (User, error) {
	var dbuser User
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&dbuser)

	return dbuser, err
}
