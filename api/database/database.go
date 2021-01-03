package database

import (
	"context"
	"github.com/shitpostingio/analysis-api/configuration/structs"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	opDeadline          = 10 * time.Second
	tokenCollectionName = "tokens"
)

var (
	dsCtx           context.Context
	database        *mongo.Database
	tokenCollection *mongo.Collection
)

// Connect connects to the MongoDB document store.
func Connect(cfg *structs.DatabaseConfiguration) {

	client, err := mongo.Connect(context.Background(), cfg.MongoDBConnectionOptions())
	if err != nil {
		log.Fatal("Unable to connect to document store:", err)
	}

	pingCtx, cancelPingCtx := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelPingCtx()
	err = client.Ping(pingCtx, readpref.Primary())
	if err != nil {
		log.Fatal("Unable to ping document store:", err)
	}

	//
	dsCtx = context.TODO()

	//
	database = client.Database(cfg.DatabaseName)
	tokenCollection = database.Collection(tokenCollectionName)

}
