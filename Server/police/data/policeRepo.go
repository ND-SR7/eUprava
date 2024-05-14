package data

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type PoliceRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// Constructor
func New(ctx context.Context, logger *log.Logger) (*PoliceRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	cr := &PoliceRepo{
		cli:    client,
		logger: logger,
	}

	return cr, nil
}

// Disconnect
func (pr *PoliceRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		pr.logger.Fatal(err.Error())
		return err
	}
	return nil
}

// Check database connection
func (pr *PoliceRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping DB
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err.Error())
	}

	// Print DBs
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err.Error())
	}
	pr.logger.Println(databases)
}

func (pr *PoliceRepo) CreateTrafficViolation(ctx context.Context, violation *TrafficViolation) error {
	collection := pr.getPoliceCollection("traffic_violations")
	_, err := collection.InsertOne(ctx, violation)
	return err
}

func (pr *PoliceRepo) GetTrafficViolationByID(ctx context.Context, id primitive.ObjectID) (*TrafficViolation, error) {
	collection := pr.getPoliceCollection("traffic_violations")
	violation := &TrafficViolation{}
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(violation)
	if err != nil {
		return nil, err
	}
	return violation, nil
}

func (pr *PoliceRepo) GetAllTrafficViolations(ctx context.Context) ([]*TrafficViolation, error) {
	collection := pr.getPoliceCollection("traffic_violations")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var violations []*TrafficViolation
	for cursor.Next(ctx) {
		var violation TrafficViolation
		if err := cursor.Decode(&violation); err != nil {
			return nil, err
		}
		violations = append(violations, &violation)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return violations, nil
}

func (pr *PoliceRepo) UpdateTrafficViolation(ctx context.Context, id primitive.ObjectID, update *TrafficViolation) error {
	collection := pr.getPoliceCollection("traffic_violations")
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (pr *PoliceRepo) DeleteTrafficViolation(ctx context.Context, id primitive.ObjectID) error {
	collection := pr.getPoliceCollection("traffic_violations")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (pr *PoliceRepo) getPoliceCollection(nameOfCollection string) *mongo.Collection {
	policeDatabase := pr.cli.Database("police_db")
	policeCollection := policeDatabase.Collection(nameOfCollection)
	return policeCollection
}
