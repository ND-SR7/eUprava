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

type MUPRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// Constructor
func New(ctx context.Context, logger *log.Logger) (*MUPRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	mr := &MUPRepo{
		cli:    client,
		logger: logger,
	}

	return mr, nil
}

// Disconnect
func (mr *MUPRepo) Disconnect(ctx context.Context) error {
	err := mr.cli.Disconnect(ctx)
	if err != nil {
		mr.logger.Fatal(err.Error())
		return err
	}
	return nil
}

// Check database connection
func (mr *MUPRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := mr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		mr.logger.Println(err)
	}

	// Print available databases
	databases, err := mr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		mr.logger.Println(err)
	}
	mr.logger.Println(databases)
}
