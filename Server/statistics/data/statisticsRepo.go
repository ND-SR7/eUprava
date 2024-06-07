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

// Initialize database
func (sr *StatisticsRepo) Initialize(ctx context.Context) error {
	db := sr.cli.Database("statisticsDB")

	err := db.Collection("trafficCollection").Drop(ctx)
	if err != nil {
		return err
	}

	trafficData := []TrafficData{
		{
			ID: primitive.NewObjectID(),
			StatisticsData: StatisticsData{
				ID:     primitive.NewObjectID(),
				Date:   time.Now(),
				Region: "Vojvodina",
				Year:   2024,
				Month:  7,
			},
			ViolationType: "Speeding",
			Vehicles: []Vehicle{
				{
					ID:           primitive.NewObjectID(),
					Brand:        "Opel",
					Model:        "Corsa C",
					Year:         2002,
					Registration: "7258",
					Plates:       "EU7258CR",
					Owner:        "147258369",
				},
			},
		},
	}

	var bsonTrafficData []interface{}
	for _, td := range trafficData {
		bsonTrafficData = append(bsonTrafficData, td)
	}

	_, err = db.Collection("trafficCollection").InsertMany(ctx, bsonTrafficData)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Repo methods

func (sr *StatisticsRepo) CreateTrafficStatisticData(ctx context.Context, trafficStatistic *TrafficData) error {
	collection := sr.getTrafficStatisticsCollection()

	_, err := collection.InsertOne(ctx, trafficStatistic)
	if err != nil {
		sr.logger.Println("Failed to create traffic statistic")
		return err
	}

	return nil
}

func (sr *StatisticsRepo) GetAllTrafficStatisticsData(ctx context.Context) ([]*TrafficData, error) {
	collection := sr.getTrafficStatisticsCollection()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		sr.logger.Println("Failed to get all traffic statistics")
		return nil, err
	}
	defer cursor.Close(ctx)

	var statistics []*TrafficData
	if err := cursor.All(ctx, &statistics); err != nil {
		sr.logger.Println("Failed to iterate over all traffic statistics")
		return nil, err
	}

	return statistics, nil
}

func (sr *StatisticsRepo) GetTrafficStatistic(ctx context.Context, id primitive.ObjectID) (*TrafficData, error) {
	collection := sr.getTrafficStatisticsCollection()

	var statistic TrafficData
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&statistic)
	if err != nil {
		sr.logger.Println("Failed to get traffic statistic")
		return nil, err
	}

	return &statistic, nil
}

func (sr *StatisticsRepo) UpdateTrafficStatistic(ctx context.Context, statistic *TrafficData) error {
	collection := sr.getTrafficStatisticsCollection()

	filter := bson.M{"_id": statistic.ID}
	update := bson.M{"$set": statistic}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		sr.logger.Println("Failed to update traffic statistic")
		return err
	}

	return nil
}

func (sr *StatisticsRepo) DeleteTrafficStatistic(ctx context.Context, id primitive.ObjectID) error {
	collection := sr.getTrafficStatisticsCollection()

	filter := bson.M{"_id": id}
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		sr.logger.Println("Failed to delete traffic statistic")
		return err
	}

	return nil
}

func (sr *StatisticsRepo) getTrafficStatisticsCollection() *mongo.Collection {
	statisticsDatabase := sr.cli.Database("statisticsDB")
	trafficStatisticsCollection := statisticsDatabase.Collection("trafficCollection")
	return trafficStatisticsCollection
}
