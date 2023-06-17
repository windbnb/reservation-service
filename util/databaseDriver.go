package util

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func ConnectToDatabase() *mongo.Database {
	connectionString, connectionStringFound := os.LookupEnv("DATABASE_CONNECTION_STRING")
	if !connectionStringFound {
		connectionString = "mongodb://user:pass@localhost:27017"
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI))

	if err != nil {
		log.Printf("Connecting to database failed.")
		return nil
	}

	return client.Database("reservation_database")
}
