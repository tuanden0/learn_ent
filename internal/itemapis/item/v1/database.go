package v1

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGO_USER_KEY = "MONGO_USERNAME"
	MONGO_PWD_KEY  = "MONGO_PWD"
)

func GetClientConnection() (*mongo.Client, error) {

	// Load mongo credentials
	if err := godotenv.Load("internal/itemapis/item/v1/.env"); err != nil {
		return nil, err
	}

	// Get mongo credentials
	usr := os.Getenv(MONGO_USER_KEY)
	pwd := os.Getenv(MONGO_PWD_KEY)

	if len(usr) == 0 || len(pwd) == 0 {
		return nil, fmt.Errorf("failed to get mongodb credentials")
	}

	// Create mongo URI
	uri := fmt.Sprintf("mongodb://%v:%v@localhost", usr, pwd)

	// Connect to mongodb
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Check client
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return client, nil
}
