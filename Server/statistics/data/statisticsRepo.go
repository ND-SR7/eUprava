package data

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type StatisticsRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// Constructor
func New(ctx context.Context, logger *log.Logger) (*StatisticsRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	sr := &StatisticsRepo{
		cli:    client,
		logger: logger,
	}

	return sr, nil
}

// Disconnect
func (sr *StatisticsRepo) Disconnect(ctx context.Context) error {
	err := sr.cli.Disconnect(ctx)
	if err != nil {
		sr.logger.Fatal(err.Error())
		return err
	}
	return nil
}

// Check database connection
func (sr *StatisticsRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := sr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		sr.logger.Println(err.Error())
	}

	// Print available databases
	databases, err := sr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		sr.logger.Println(err.Error())
	}
	sr.logger.Println(databases)
}

// TODO: Repo methods
