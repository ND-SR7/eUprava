package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

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

// Initialize database
func (mr *MUPRepo) Initialize(ctx context.Context) error {
	db := mr.cli.Database("mupDB")

	err := db.Collection("vehicle").Drop(ctx)
	if err != nil {
		return err
	}

	vehicles := []Vehicle{
		{
			ID:           primitive.NewObjectID(),
			Brand:        "Toyota",
			Model:        "Corolla",
			Year:         2020,
			Registration: "",
			Plates:       "",
			Owner:        "123456789",
		},
		{
			ID:           primitive.NewObjectID(),
			Brand:        "Honda",
			Model:        "Civic",
			Year:         2019,
			Registration: "",
			Plates:       "",
			Owner:        "123456789",
		},
		{
			ID:           primitive.NewObjectID(),
			Brand:        "Ford",
			Model:        "Focus",
			Year:         2018,
			Registration: "",
			Plates:       "",
			Owner:        "33355577799",
		},
		{
			ID:           primitive.NewObjectID(),
			Brand:        "Chevrolet",
			Model:        "Malibu",
			Year:         2021,
			Registration: "",
			Plates:       "",
			Owner:        "33355577799",
		},
	}

	var bsonVehicles []interface{}
	for _, v := range vehicles {
		bsonVehicles = append(bsonVehicles, v)
	}

	_, err = db.Collection("vehicle").InsertMany(ctx, bsonVehicles)
	if err != nil {
		return err
	}

	return nil
}

//Vehicle methods

func (mr *MUPRepo) SaveVehicle(ctx context.Context, vehicle *Vehicle) error {
	collection := mr.getMupCollection("vehicle")
	println("owner ", vehicle.Owner)

	_, err := collection.InsertOne(ctx, vehicle)
	if err != nil {
		log.Printf("Failed to create vehicle: %v", err)
		return err
	}

	err = mr.SaveVehicleIntoMup(ctx, vehicle)
	if err != nil {
		log.Printf("Failed to save vehicle into mup: %v", err)
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

func (mr *MUPRepo) GetPlatesByVehicleID(ctx context.Context, vehicleID primitive.ObjectID) (Plates, error) {
	collection := mr.getMupCollection("plates")
	var plates Plates
	err := collection.FindOne(ctx, bson.M{"vehicleID": vehicleID}).Decode(&plates)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Plates{}, nil
		}
		return Plates{}, err
	}
	return plates, nil
}

func (mr *MUPRepo) GetRegistrationByVehicleID(ctx context.Context, vehicleID primitive.ObjectID) (Registration, error) {
	collection := mr.getMupCollection("registration")
	var registration Registration
	err := collection.FindOne(ctx, bson.M{"vehicleID": vehicleID}).Decode(&registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Registration{}, nil
		}
		return Registration{}, err
	}
	return registration, nil
}

//Mup methods

func (mr *MUPRepo) SaveVehicleIntoMup(ctx context.Context, vehicle *Vehicle) error {
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"name", "Mup"}}

	update := bson.D{{"$push", bson.D{{"vehicles", vehicle.ID}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Vehicle successfully saved into mup!")
	return nil
}

func (mr *MUPRepo) SavePlatesIntoMup(ctx context.Context, plates Plates) error {
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"name", "Mup"}}

	update := bson.D{{"$push", bson.D{{"plates", plates.PlatesNumber}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Plates successfully saved into mup!")
	return nil
}

func (mr *MUPRepo) SaveRegistrationIntoMup(ctx context.Context, registration *Registration) error {
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"name", "Mup"}}

	update := bson.D{{"$push", bson.D{{"registrations", registration.RegistrationNumber}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Registration successfully saved into mup!")
	return nil
}

func (mr *MUPRepo) SaveTrafficPermitIntoMup(ctx context.Context, trafficPermit *TrafficPermit) error {
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"name", "Mup"}}

	update := bson.D{{"$push", bson.D{{"trafficPermits", trafficPermit.ID}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Traffic permit successfully saved into mup!")
	return nil
}

func (mr *MUPRepo) SaveDrivingBanIntoMup(ctx context.Context, drivingBan DrivingBan) error {
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"name", "Mup"}}

	update := bson.D{{"$push", bson.D{{"drivingBans", drivingBan.ID}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println("Driving ban successfully saved into mup!")
	return nil
}

func (mr *MUPRepo) SubmitRegistrationRequest(ctx context.Context, registration *Registration) error {
	collection := mr.getMupCollection("registration")

	_, err := collection.InsertOne(ctx, registration)
	if err != nil {
		log.Printf("Failed to create vehicle: %v", err)
		return err
	}

	err = mr.SaveRegistrationIntoMup(ctx, registration)
	if err != nil {
		log.Printf("Failed to save registration into mup: %v", err)
		return err
	}

	err = mr.SaveRegistrationIntoVehicle(ctx, registration)
	if err != nil {
		log.Printf("Failed to save registration into vehicle: %v", err)
		return err
	}

	return nil
}

func (mr *MUPRepo) ApproveRegistration(ctx context.Context, registration Registration) error {
	collection := mr.getMupCollection("registration")

	filter := bson.D{{"registrationNumber", registration.RegistrationNumber}}

	update := bson.D{{"$set", bson.D{
		{"approved", true},
		{"expirationDate", registration.ExpirationDate},
		{"plates", registration.Plates}}}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {

		return err
	}

	fmt.Println("Traffic permit approved successfully!")

	return nil
}

func (mr *MUPRepo) DeletePendingRegistration(ctx context.Context, registrationNumber string) error {
	collection := mr.getMupCollection("registration")
	filter := bson.D{{"registrationNumber", registrationNumber}, {"approved", false}}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Println("Pending registration request deleted successfully!")
	return nil
}

// Delete pending traffic permit request
func (mr *MUPRepo) DeletePendingTrafficPermit(ctx context.Context, permitID primitive.ObjectID) error {
	collection := mr.getMupCollection("trafficPermit")
	filter := bson.D{{"_id", permitID}, {"approved", false}}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	fmt.Println("Pending traffic permit request deleted successfully!")
	return nil
}

func (mr *MUPRepo) GetRegistrationByPlate(ctx context.Context, plate string) (Registration, error) {
	collection := mr.getMupCollection("registration")

	filter := bson.D{{"plates", plate}}

	var registration Registration

	err := collection.FindOne(ctx, filter).Decode(&registration)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Registration{}, nil
		}
		return Registration{}, err
	}

	return registration, nil
}

//Driving permit methods

func (mr *MUPRepo) SubmitTrafficPermitRequest(ctx context.Context, trafficPermit *TrafficPermit) error {
	collection := mr.getMupCollection("trafficPermit")

	_, err := collection.InsertOne(ctx, trafficPermit)
	if err != nil {
		log.Printf("Failed to create traffic permit: %v", err)
		return err
	}

	err = mr.SaveTrafficPermitIntoMup(ctx, trafficPermit)
	if err != nil {
		log.Printf("Failed to save traffic permit into mup: %v", err)
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
		log.Printf("Failed to create plates: %v", err)
		return err
	}

	err = mr.SavePlatesIntoMup(ctx, plates)
	if err != nil {
		log.Printf("Failed to save plates into mup: %v", err)
		return err
	}

	err = mr.SavePlatesIntoVehicle(ctx, plates)
	if err != nil {
		log.Printf("Failed to save plates into vehicle: %v", err)
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
		log.Printf("Failed to create driving ban: %v", err)
		return err
	}

	err = mr.SaveDrivingBanIntoMup(ctx, *drivingBan)
	if err != nil {
		log.Printf("Failed to save driving ban into mup: %v", err)
		return err
	}

	return nil
}

func (mr *MUPRepo) GetDrivingBan(ctx context.Context, jmbg string, now time.Time) (DrivingBan, error) {
	collection := mr.getMupCollection("drivingBan")

	filter := bson.D{
		{"person", jmbg},
		{"duration", bson.D{{"$lt", now}}},
	}

	var drivingBan DrivingBan

	err := collection.FindOne(ctx, filter).Decode(&drivingBan)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return DrivingBan{}, nil
		}
		return DrivingBan{}, err
	}

	return drivingBan, nil
}

func (mr *MUPRepo) GetDrivingPermitByJMBG(ctx context.Context, jmbg string) (TrafficPermit, error) {
	collection := mr.getMupCollection("trafficPermit")

	filter := bson.D{{"person", jmbg}}

	var drivingPermit TrafficPermit

	err := collection.FindOne(ctx, filter).Decode(&drivingPermit)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return TrafficPermit{}, nil
		}
		return TrafficPermit{}, err
	}

	return drivingPermit, nil
}

//Person methods

func (mr *MUPRepo) CheckForPersonsDrivingBans(ctx context.Context, userID string) (DrivingBans, error) {
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

func (mr *MUPRepo) GetPersonsRegistrations(ctx context.Context, jmbg string) (Registrations, error) {
	collection := mr.getMupCollection("registration")

	filter := bson.D{{"owner", jmbg}}

	var registrations Registrations

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var registration Registration
		if err := cursor.Decode(&registration); err != nil {
			return nil, err
		}
		registrations = append(registrations, registration)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return registrations, nil
}

func (mr *MUPRepo) GetUserDrivingPermit(ctx context.Context, jmbg string) (TrafficPermits, error) {
	collection := mr.getMupCollection("trafficPermit")

	filter := bson.D{{"person", jmbg}, {"approved", true}}

	var drivingPermits TrafficPermits

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var drivingPermit TrafficPermit
		if err := cursor.Decode(&drivingPermit); err != nil {
			return nil, err
		}
		drivingPermits = append(drivingPermits, drivingPermit)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return drivingPermits, nil
}

func (mr *MUPRepo) GetPersonsVehicles(ctx context.Context, jmbg string) ([]Vehicle, error) {
	collection := mr.getMupCollection("vehicle")

	filter := bson.M{"owner": jmbg}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Failed to find vehicles for owner %s: %v", jmbg, err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var vehicles []Vehicle
	if err = cursor.All(ctx, &vehicles); err != nil {
		log.Printf("Failed to decode vehicles: %v", err)
		return nil, err
	}

	return vehicles, nil
}

func (mr *MUPRepo) GetPendingRegistrationRequests(ctx context.Context) (Registrations, error) {
	collection := mr.getMupCollection("registration")

	filter := bson.D{
		{"approved", false},
	}

	var pendingRequests Registrations

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var registration Registration
		if err := cursor.Decode(&registration); err != nil {
			return nil, err
		}
		pendingRequests = append(pendingRequests, registration)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return pendingRequests, nil
}

func (mr *MUPRepo) GetPendingTrafficPermitRequests(ctx context.Context) (TrafficPermits, error) {
	collection := mr.getMupCollection("trafficPermit")

	filter := bson.D{
		{"approved", false},
	}

	var pendingRequests TrafficPermits

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var trafficPermit TrafficPermit
		if err := cursor.Decode(&trafficPermit); err != nil {
			return nil, err
		}
		pendingRequests = append(pendingRequests, trafficPermit)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return pendingRequests, nil
}

// MUP methods
func (mr *MUPRepo) SaveMup(ctx context.Context) error {
	collection := mr.getMupCollection("mup")

	filter := bson.D{{"name", "Mup"}}
	var existingMup Mup
	err := collection.FindOne(ctx, filter).Decode(&existingMup)
	if err == nil {
		log.Printf("Mup with name 'Mup' already exists: %v", existingMup)
		return nil
	}
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Error while checking if Mup exists: %v", err)
		return err
	}

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

	_, err = collection.InsertOne(ctx, mup)
	if err != nil {
		log.Printf("Failed to create Mup: %v", err)
		return err
	}

	log.Printf("Inserted Mup: %v", mup)
	return nil
}

// Get collection method
func (mr *MUPRepo) getMupCollection(nameOfCollection string) *mongo.Collection {
	mupDatabase := mr.cli.Database("mupDB")
	mupCollection := mupDatabase.Collection(nameOfCollection)
	return mupCollection
}
