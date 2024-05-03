package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
	//dburi := os.Getenv("mongodb://root:pass@mup_db:27017")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:pass@mup_db:27017"))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create collection
	collection := client.Database("mup_db").Collection("person")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(collection.Name())

	fmt.Println("Connected to db")

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

//Vehicle methods

func (mr *MUPRepo) SaveVehicle(ctx context.Context, vehicle Vehicle) error {
	collection := mr.getMupCollection("vehicle")

	_, err := collection.InsertOne(ctx, vehicle)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create vehicle: %v", err))
		return err
	}

	return nil
}

func (mr *MUPRepo) SaveRegistrationIntoVehicle(ctx context.Context, registration Registration) error {
	collection := mr.getMupCollection("vehicle")

	filter := bson.D{{"_id", registration.VehicleID}}

	update := bson.D{{"$set", bson.D{{"registration", registration.RegistrationNumber}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Registration successfully saved into vehicle!")
	return nil
}

func (mr *MUPRepo) SavePlatesIntoVehicle(ctx context.Context, plates Plates) error {
	collection := mr.getMupCollection("vehicle")

	filter := bson.D{{"_id", plates.VehicleID}}

	update := bson.D{{"$set", bson.D{{"plates", plates.PlatesNumber}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Plates successfully saved into vehicle!")
	return nil
}

//Mup methods

func (mr *MUPRepo) SaveMup(ctx context.Context, mup Mup) error {
	collection := mr.getMupCollection("mup")

	_, err := collection.InsertOne(ctx, mup)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create mup: %v", err))
		return err
	}

	return nil
}

func (mr *MUPRepo) SavePlatesIntoMup(ctx context.Context, plates Plates) error {
	mupID, err := primitive.ObjectIDFromHex("607d22b837ede6b71eef3e82")
	if err != nil {
		return err
	}
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"_id", mupID}}

	update := bson.D{{"$push", bson.D{{"plates", plates.PlatesNumber}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Plates successfully saved into mup!")
	return nil
}

func (mr *MUPRepo) SaveRegistrationIntoMup(ctx context.Context, registration Registration) error {
	mupID, err := primitive.ObjectIDFromHex("607d22b837ede6b71eef3e82")
	if err != nil {
		return err
	}
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"_id", mupID}}

	update := bson.D{{"$push", bson.D{{"registrations", registration.RegistrationNumber}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Registration successfully saved into mup!")
	return nil
}

func (mr *MUPRepo) SaveTrafficPermitIntoMup(ctx context.Context, trafficPermit TrafficPermit) error {
	mupID, err := primitive.ObjectIDFromHex("607d22b837ede6b71eef3e82")
	if err != nil {
		return err
	}
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"_id", mupID}}

	update := bson.D{{"$push", bson.D{{"trafficPermits", trafficPermit.ID}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Traffic permit successfully saved into mup!")
	return nil
}

func (mr *MUPRepo) SaveDrivingBanIntoMup(ctx context.Context, drivingBan DrivingBan) error {
	mupID, err := primitive.ObjectIDFromHex("607d22b837ede6b71eef3e82")
	if err != nil {
		return err
	}
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"_id", mupID}}

	update := bson.D{{"$push", bson.D{{"drivingBans", drivingBan.ID}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Driving ban successfully saved into mup!")
	return nil
}

//Registration methods

//func (mr *MUPRepo) RegisterVehicle(ctx context.Context, registration Registration) error {
//	registration.IssuedDate = time.Now()
//	registration.Approved = false
//
//	collection := mr.getMupCollection("registration")
//
//	_, err := collection.InsertOne(ctx, registration)
//	if err != nil {
//		log.Printf(fmt.Sprintf("Failed to create registration: %v", err))
//		return err
//	}
//
//	err = mr.SaveRegistrationIntoMup(ctx, registration)
//	if err != nil {
//		log.Printf(fmt.Sprintf("Failed to save registration into mup: %v", err))
//		return err
//	}
//
//	err = mr.SaveRegistrationIntoVehicle(ctx, registration)
//	if err != nil {
//		log.Printf(fmt.Sprintf("Failed to save registration into vehicle: %v", err))
//		return err
//	}
//
//	return nil
//}

func (mr *MUPRepo) SubmitRegistrationRequest(ctx context.Context, registration Registration) error {
	registration.Approved = false
	registration.IssuedDate = time.Now()
	collection := mr.getMupCollection("registration")

	_, err := collection.InsertOne(ctx, registration)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create vehicle: %v", err))
		return err
	}

	return nil
}

func (mr *MUPRepo) ApproveRegistration(ctx context.Context, registrationID primitive.ObjectID) error {
	expirationDate := time.Now().AddDate(1, 0, 0)

	collection := mr.getMupCollection("registration")

	filter := bson.D{{"_id", registrationID}}

	update := bson.D{{"$set", bson.D{{"approved", true}, {"expirationDate", expirationDate}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Traffic permit approved successfully!")
	return nil
}

//Driving permit methods

func (mr *MUPRepo) SubmitTrafficPermitRequest(ctx context.Context, trafficPermit TrafficPermit) error {
	trafficPermit.Approved = false
	trafficPermit.IssuedDate = time.Now()

	collection := mr.getMupCollection("trafficPermit")

	_, err := collection.InsertOne(ctx, trafficPermit)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create traffic permit: %v", err))
		return err
	}

	err = mr.SaveTrafficPermitIntoMup(ctx, trafficPermit)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to save traffic permit into mup: %v", err))
		return err
	}

	return nil
}

func (mr *MUPRepo) ApproveTrafficPermitRequest(ctx context.Context, permitID primitive.ObjectID) error {
	expirationDate := time.Now().AddDate(5, 0, 0)

	collection := mr.getMupCollection("trafficPermit")

	filter := bson.D{{"_id", permitID}}

	update := bson.D{{"$set", bson.D{{"approved", true}, {"expirationDate", expirationDate}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Traffic permit approved successfully!")
	return nil
}

//Plates methods

func (mr *MUPRepo) IssuePlates(ctx context.Context, plates Plates) error {
	collection := mr.getMupCollection("plates")

	_, err := collection.InsertOne(ctx, plates)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create plates: %v", err))
		return err
	}

	err = mr.SavePlatesIntoMup(ctx, plates)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to save plates into mup: %v", err))
		return err
	}

	err = mr.SavePlatesIntoVehicle(ctx, plates)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to save plates into vehicle: %v", err))
		return err
	}

	return nil
}

//Driving ban methods

func (mr *MUPRepo) IssueDrivingBan(ctx context.Context, drivingBan DrivingBan) error {
	collection := mr.getMupCollection("drivingBan")

	_, err := collection.InsertOne(ctx, drivingBan)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create driving ban: %v", err))
		return err
	}

	err = mr.SaveDrivingBanIntoMup(ctx, drivingBan)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to save driving ban into mup: %v", err))
		return err
	}

	return nil
}

func (mr *MUPRepo) CheckForDrivingBan(userID primitive.ObjectID) ([]DrivingBan, error) {
	collection := mr.getMupCollection("drivingBan")

	filter := bson.D{{"person", userID}}

	var drivingBans []DrivingBan

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var drivingBan DrivingBan
		if err := cursor.Decode(&drivingBan); err != nil {
			return nil, err
		}
		drivingBans = append(drivingBans, drivingBan)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return drivingBans, nil
}

//Get collection method

func (mr *MUPRepo) getMupCollection(nameOfCollection string) *mongo.Collection {
	mupDatabase := mr.cli.Database("mup_db")
	mupCollection := mupDatabase.Collection(nameOfCollection)
	return mupCollection
}
