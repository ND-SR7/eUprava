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

type MupRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// Constructor
func New(ctx context.Context, logger *log.Logger) (*MupRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	nr := &MupRepo{
		cli:    client,
		logger: logger,
	}

	return nr, nil
}

// Disconnect
func (sr *MupRepo) Disconnect(ctx context.Context) error {
	err := sr.cli.Disconnect(ctx)
	if err != nil {
		sr.logger.Fatal(err.Error())
		return err
	}
	return nil
}

// Check database connection
func (sr *MupRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := sr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		sr.logger.Println(err)
	}

	// Print available databases
	databases, err := sr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		sr.logger.Println(err)
	}
	sr.logger.Println(databases)
}
