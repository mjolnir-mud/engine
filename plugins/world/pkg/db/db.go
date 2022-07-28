package db

import (
	"context"
	"os"

	"github.com/mjolnir-mud/engine/plugins/world/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var log = logger.Logger.With().Str("service", "db").Logger()
var client *mongo.Client

func Start() {
	log.Info().Msg("starting db connection")

	c, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to mongo")
		os.Exit(1)
	}

	err = c.Connect(context.Background())

	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to mongo")
		os.Exit(1)
	}

	client = c
}

func Collection(name string) *mongo.Collection {
	return client.Database("world").Collection(name)
}
