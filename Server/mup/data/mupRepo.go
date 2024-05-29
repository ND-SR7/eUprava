package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	cr := &MUPRepo{
		cli:    client,
		logger: logger,
	}

	return cr, nil
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

func (mr *MUPRepo) SaveVehicle(ctx context.Context, vehicle *Vehicle) error {
	vehicle.ID = primitive.NewObjectID()
	vehicle.Registration = ""
	vehicle.Plates = ""
	collection := mr.getMupCollection("vehicle")

	_, err := collection.InsertOne(ctx, vehicle)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create vehicle: %v", err))
		return err
	}

	err = mr.SaveVehicleIntoMup(ctx, vehicle)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to save vehicle into mup: %v", err))
		return err
	}

	return nil
}

func (mr *MUPRepo) SaveRegistrationIntoVehicle(ctx context.Context, registration *Registration) error {
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

func (mr *MUPRepo) RetrieveRegisteredVehicles(ctx context.Context) (Vehicles, error) {
	collection := mr.getMupCollection("vehicle")

	filter := bson.M{"registration": bson.M{"$ne": ""}}

	var vehicles Vehicles

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var vehicle Vehicle
		if err := cursor.Decode(&vehicle); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, vehicle)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return vehicles, nil
}

//Mup methods

func (mr *MUPRepo) SaveVehicleIntoMup(ctx context.Context, vehicle *Vehicle) error {
	mupID, err := primitive.ObjectIDFromHex("607d22b837ede6b71eef3e82")
	if err != nil {
		return err
	}
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"_id", mupID}}

	update := bson.D{{"$push", bson.D{{"vehicles", vehicle.ID}}}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Vehicle successfully saved into mup!")
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

func (mr *MUPRepo) SaveRegistrationIntoMup(ctx context.Context, registration *Registration) error {
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

func (mr *MUPRepo) SaveTrafficPermitIntoMup(ctx context.Context, trafficPermit *TrafficPermit) error {
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

func (mr *MUPRepo) SubmitRegistrationRequest(ctx context.Context, registration *Registration) error {
	registration.Approved = false
	registration.IssuedDate = time.Now()
	collection := mr.getMupCollection("registration")

	_, err := collection.InsertOne(ctx, registration)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create vehicle: %v", err))
		return err
	}

	err = mr.SaveRegistrationIntoMup(ctx, registration)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to save registration into mup: %v", err))
		return err
	}

	err = mr.SaveRegistrationIntoVehicle(ctx, registration)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to save registration into vehicle: %v", err))
		return err
	}

	return nil
}

func (mr *MUPRepo) ApproveRegistration(ctx context.Context, registration Registration) error {
	expirationDate := time.Now().AddDate(1, 0, 0)

	registration.Approved = true

	collection := mr.getMupCollection("registration")

	filter := bson.D{{"_id", registration.RegistrationNumber}}

	update := bson.D{{"$set", bson.D{{"approved", true}, {"expirationDate", expirationDate}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {

		return err
	}

	fmt.Println("Traffic permit approved successfully!")

	plates := Plates{
		RegistrationNumber: registration.RegistrationNumber,
		PlatesNumber:       "222222",
		PlateType:          "aaa",
		VehicleID:          registration.VehicleID,
	}

	err = mr.IssuePlates(ctx, plates)
	if err != nil {
		return err
	}

	return nil
}

//Driving permit methods

func (mr *MUPRepo) SubmitTrafficPermitRequest(ctx context.Context, trafficPermit *TrafficPermit) error {
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

func (mr *MUPRepo) IssueDrivingBan(ctx context.Context, drivingBan *DrivingBan) error {
	drivingBan.ID = primitive.NewObjectID()

	collection := mr.getMupCollection("drivingBan")

	_, err := collection.InsertOne(ctx, drivingBan)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to create driving ban: %v", err))
		return err
	}

	err = mr.SaveDrivingBanIntoMup(ctx, *drivingBan)
	if err != nil {
		log.Printf(fmt.Sprintf("Failed to save driving ban into mup: %v", err))
		return err
	}

	return nil
}

//Person methods

func (mr *MUPRepo) CheckForPersonsDrivingBans(ctx context.Context, userID primitive.ObjectID) (DrivingBans, error) {
	collection := mr.getMupCollection("drivingBan")

	filter := bson.D{{"person", userID}}

	var drivingBans DrivingBans

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
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

// MUP methods
func (mr *MUPRepo) SaveMup(ctx context.Context) error {
	collection := mr.getMupCollection("mup")

	mupID := primitive.NewObjectID()

	address := Address{
		Municipality: "",
		Locality:     "Novi Sad",
		StreetName:   "Dunavska",
		StreetNumber: 1,
	}

	mup := Mup{
		ID:             mupID,
		Name:           "Mup",
		Address:        address,
		Vehicles:       make([]primitive.ObjectID, 0),
		TrafficPermits: make([]primitive.ObjectID, 0),
		Plates:         make([]string, 0),
		DrivingBans:    make([]primitive.ObjectID, 0),
		Registrations:  make([]string, 0),
	}

	_, err := collection.InsertOne(ctx, mup)
	if err != nil {
		log.Printf("Failed to create mup: %v", err)
		return err
	}

	log.Printf("Inserted Mup: %v", mup)
	return nil
}

//Get collection method

func (mr *MUPRepo) getMupCollection(nameOfCollection string) *mongo.Collection {
	mupDatabase := mr.cli.Database("mup_db")
	mupCollection := mupDatabase.Collection(nameOfCollection)
	return mupCollection
}
