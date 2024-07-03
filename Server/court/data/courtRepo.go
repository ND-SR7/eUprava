package data

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type CourtRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// Constructor
func New(ctx context.Context, logger *log.Logger) (*CourtRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	cr := &CourtRepo{
		cli:    client,
		logger: logger,
	}

	return cr, nil
}

// Disconnect
func (cr *CourtRepo) Disconnect(ctx context.Context) error {
	err := cr.cli.Disconnect(ctx)
	if err != nil {
		cr.logger.Fatal(err.Error())
		return err
	}
	return nil
}

// Check database connection
func (cr *CourtRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping DB
	err := cr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		cr.logger.Println(err.Error())
	}

	// Print DBs
	databases, err := cr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		cr.logger.Println(err.Error())
	}
	cr.logger.Println(databases)
}

// Initialize database
func (cr *CourtRepo) Initialize(ctx context.Context) error {
	courtCollection := cr.getCourtCollection()
	hearingsPersonCollection := cr.getHearingsPersonCollection()
	hearingsLegalEntityCollection := cr.getHearingsLegalEntityCollection()
	suspensionsCollection := cr.getSuspensionsCollection()
	warrantsCollection := cr.getWarrantsCollection()

	err := courtCollection.Drop(ctx)
	if err != nil {
		return err
	}

	err = hearingsPersonCollection.Drop(ctx)
	if err != nil {
		return err
	}
	err = hearingsLegalEntityCollection.Drop(ctx)
	if err != nil {
		return err
	}
	err = warrantsCollection.Drop(ctx)
	if err != nil {
		return err
	}
	err = suspensionsCollection.Drop(ctx)
	if err != nil {
		return err
	}

	hearingDateTime, _ := time.Parse("2006-01-02T15:04:05", "2024-07-05T16:00:00")
	courtID, _ := primitive.ObjectIDFromHex("64c13ab08edf48a008793cac")

	court := Court{
		ID:   courtID,
		Name: "Sud eUprave",
		Address: Address{
			Municipality: "Novi Sad",
			Locality:     "Novi Sad",
			StreetName:   "Jovana Zmaja",
			StreetNumber: 3,
		},
	}

	_, err = courtCollection.InsertOne(ctx, court)
	if err != nil {
		return err
	}

	courtHearingsPerson := []CourtHearingPerson{
		{
			ID:       primitive.NewObjectID(),
			Reason:   "Speeding violation hearing",
			DateTime: hearingDateTime,
			Court:    courtID,
			Person:   "147258369",
		},
		{
			ID:       primitive.NewObjectID(),
			Reason:   "Drunk driving violation hearing",
			DateTime: hearingDateTime.Add(time.Hour),
			Court:    courtID,
			Person:   "369258147",
		},
	}

	var bsonCHP []interface{}
	for _, chp := range courtHearingsPerson {
		bsonCHP = append(bsonCHP, chp)
	}

	_, err = hearingsPersonCollection.InsertMany(ctx, bsonCHP)
	if err != nil {
		return err
	}

	courtHearingsLegalEntity := []CourtHearingLegalEntity{
		{
			ID:          primitive.NewObjectID(),
			Reason:      "Speeding violation hearing",
			DateTime:    hearingDateTime.Add(time.Hour * 2),
			Court:       courtID,
			LegalEntity: "147369258",
		},
	}

	var bsonCHLE []interface{}
	for _, chle := range courtHearingsLegalEntity {
		bsonCHLE = append(bsonCHLE, chle)
	}

	_, err = hearingsLegalEntityCollection.InsertMany(ctx, bsonCHLE)
	if err != nil {
		return err
	}

	return nil
}

// Finds court based on provided id
func (cr *CourtRepo) GetCourtByID(id string) (Court, error) {
	collection := cr.getCourtCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Court{}, err
	}

	filter := bson.M{"_id": objID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var court Court

	err = collection.FindOne(ctx, filter).Decode(&court)
	if err != nil {
		return Court{}, err
	}

	return court, nil
}

// Finds court hearing based on provided id
func (cr *CourtRepo) GetHearingByID(id string) (CourtHearing, error) {
	hearingsPerson := cr.getHearingsPersonCollection()
	hearingsLegalEntity := cr.getHearingsLegalEntityCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var courtHearing CourtHearing
	var courtHearingPerson CourtHearingPerson
	var courtHearingLegalEntity CourtHearingLegalEntity

	err = hearingsPerson.FindOne(ctx, filter).Decode(&courtHearingPerson)
	if err == nil {
		courtHearing = &courtHearingPerson
		return courtHearing, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	err = hearingsLegalEntity.FindOne(ctx, filter).Decode(&courtHearingLegalEntity)
	if err == nil {
		courtHearing = &courtHearingLegalEntity
		return courtHearing, nil
	} else if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	return nil, errors.New("hearing not found")
}

// Finds court hearings based on provided JMBG
func (cr *CourtRepo) GetHearingsByJMBG(jmbg string) ([]CourtHearing, error) {
	hearingsPerson := cr.getHearingsPersonCollection()
	hearingsLegalEntity := cr.getHearingsLegalEntityCollection()

	filterPerson := bson.M{"person": jmbg}
	filterLegalEntity := bson.M{"legalEntity": jmbg}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var courtHearings []CourtHearing

	cursorPersons, err := hearingsPerson.Find(ctx, filterPerson)
	if err != nil {
		return nil, err
	}

	for cursorPersons.Next(ctx) {
		var chp CourtHearingPerson
		if err := cursorPersons.Decode(&chp); err != nil {
			return nil, err
		}
		courtHearings = append(courtHearings, &chp)
	}
	if err := cursorPersons.Err(); err != nil {
		return nil, err
	}

	cursorLegalEntities, err := hearingsLegalEntity.Find(ctx, filterLegalEntity)
	if err != nil {
		return nil, err
	}

	for cursorLegalEntities.Next(ctx) {
		var chle CourtHearingLegalEntity
		if err := cursorPersons.Decode(&chle); err != nil {
			return nil, err
		}
		courtHearings = append(courtHearings, &chle)
	}
	if err := cursorLegalEntities.Err(); err != nil {
		return nil, err
	}

	if len(courtHearings) > 0 {
		return courtHearings, nil
	}

	return nil, errors.New("hearings not found")
}

// Finds suspension based on provided JMBG
func (cr *CourtRepo) GetSuspensionByJMBG(jmbg string) (Suspension, error) {
	collection := cr.getSuspensionsCollection()

	filter := bson.M{"person": jmbg}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var suspension Suspension

	err := collection.FindOne(ctx, filter).Decode(&suspension)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return Suspension{}, err
	}

	return suspension, nil
}

// Finds warrants based on provided JMBG
func (cr *CourtRepo) GetWarrantsByJMBG(jmbg string) (Warrants, error) {
	collection := cr.getWarrantsCollection()

	filter := bson.M{"issuedFor": jmbg}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var warrants Warrants
	for cursor.Next(ctx) {
		var warrant Warrant
		if err := cursor.Decode(&warrant); err != nil {
			return nil, err
		}
		warrants = append(warrants, warrant)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return warrants, nil
}

// Inserts new court hearing for a person into collection
func (cr *CourtRepo) CreateHearingPerson(newHearing NewCourtHearingPerson) error {
	collection := cr.getHearingsPersonCollection()

	dateTime, err := time.Parse("2006-01-02T15:04:05", newHearing.DateTime)
	if err != nil {
		cr.logger.Println("Error while parsing date")
		return err
	}

	courtID, err := primitive.ObjectIDFromHex(newHearing.Court)
	if err != nil {
		cr.logger.Println("Error while parsing court ID")
		return err
	}

	hearing := CourtHearingPerson{
		ID:       primitive.NewObjectID(),
		Reason:   newHearing.Reason,
		DateTime: dateTime,
		Court:    courtID,
		Person:   newHearing.Person,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, hearing)
	if err != nil {
		cr.logger.Fatalln("Failed to insert new court hearing for person")
		return err
	}

	cr.logger.Println(hearing)
	return nil
}

// Inserts new court hearing for a legal entity into collection
func (cr *CourtRepo) CreateHearingLegalEntity(newHearing NewCourtHearingLegalEntity) error {
	collection := cr.getHearingsLegalEntityCollection()

	dateTime, err := time.Parse("2006-01-02T15:04:05", newHearing.DateTime)
	if err != nil {
		cr.logger.Println("Error while parsing date")
		return err
	}

	courtID, err := primitive.ObjectIDFromHex(newHearing.Court)
	if err != nil {
		cr.logger.Println("Error while parsing court ID")
		return err
	}

	hearing := CourtHearingLegalEntity{
		ID:          primitive.NewObjectID(),
		Reason:      newHearing.Reason,
		DateTime:    dateTime,
		Court:       courtID,
		LegalEntity: newHearing.LegalEntity,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, hearing)
	if err != nil {
		cr.logger.Fatalln("Failed to insert new court hearing for legal entity")
		return err
	}

	cr.logger.Println(hearing)
	return nil
}

// Inserts a new warrant into collection
func (cr *CourtRepo) CreateWarrant(newWarrant NewWarrant) error {
	collection := cr.getWarrantsCollection()

	trafficViolationID, err := primitive.ObjectIDFromHex(newWarrant.TrafficViolation)
	if err != nil {
		cr.logger.Println("Error while parsing traffic violation ID")
		return err
	}

	warrant := Warrant{
		ID:               primitive.NewObjectID(),
		TrafficViolation: trafficViolationID,
		IssuedOn:         time.Now(),
		IssuedFor:        newWarrant.IssuedFor,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, warrant)
	if err != nil {
		cr.logger.Fatalln("Failed to insert new warrant")
		return err
	}

	return nil
}

// Inserts a new suspension into collection
func (cr *CourtRepo) CreateSuspension(newSuspension NewSuspension) error {
	collection := cr.getSuspensionsCollection()

	fromDateTime, err := time.Parse("2006-01-02T15:04:05", newSuspension.From)
	if err != nil {
		cr.logger.Println("Error while parsing date")
		return err
	}

	toDateTime, err := time.Parse("2006-01-02T15:04:05", newSuspension.To)
	if err != nil {
		cr.logger.Println("Error while parsing date")
		return err
	}

	suspension := Suspension{
		ID:     primitive.NewObjectID(),
		From:   fromDateTime,
		To:     toDateTime,
		Person: newSuspension.Person,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, suspension)
	if err != nil {
		cr.logger.Fatalln("Failed to insert new suspension")
		return err
	}

	return nil
}

// Reschedules court hearing for a later date and time
func (cr *CourtRepo) RescheduleCourtHearingPerson(rescheduledHearing RescheduleCourtHearing) error {
	collection := cr.getHearingsPersonCollection()

	objID, err := primitive.ObjectIDFromHex(rescheduledHearing.HearingID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}

	dateTime, err := time.Parse("2006-01-02T15:04:05", rescheduledHearing.DateTime)
	if err != nil {
		cr.logger.Println("Error while parsing date")
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"dateTime": dateTime,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		cr.logger.Println("Error while updating collection")
		return err
	}

	if result.ModifiedCount > 0 {
		cr.logger.Println("Successfully updated collection")
		return nil
	}

	cr.logger.Printf("Invalid court hearing ID: %s", rescheduledHearing.HearingID)
	return errors.New("invalid hearing ID")
}

// Reschedules court hearing for a later date and time
func (cr *CourtRepo) RescheduleCourtHearingLegalEntity(rescheduledHearing RescheduleCourtHearing) error {
	collection := cr.getHearingsLegalEntityCollection()

	objID, err := primitive.ObjectIDFromHex(rescheduledHearing.HearingID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}

	dateTime, err := time.Parse("2006-01-02T15:04:05", rescheduledHearing.DateTime)
	if err != nil {
		cr.logger.Println("Error while parsing date")
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"dateTime": dateTime,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		cr.logger.Println("Error while updating collection")
		return err
	}

	if result.ModifiedCount > 0 {
		cr.logger.Println("Successfully updated collection")
		return nil
	}

	cr.logger.Printf("Invalid court hearing ID: %s", rescheduledHearing.HearingID)
	return errors.New("invalid hearing ID")
}

// Getters for collections

func (cr *CourtRepo) getCourtCollection() *mongo.Collection {
	return cr.cli.Database("courtDB").Collection("court")
}

func (cr *CourtRepo) getHearingsPersonCollection() *mongo.Collection {
	return cr.cli.Database("courtDB").Collection("hearingsPerson")
}

func (cr *CourtRepo) getHearingsLegalEntityCollection() *mongo.Collection {
	return cr.cli.Database("courtDB").Collection("hearingsLegalEntity")
}

func (cr *CourtRepo) getWarrantsCollection() *mongo.Collection {
	return cr.cli.Database("courtDB").Collection("warrants")
}

func (cr *CourtRepo) getSuspensionsCollection() *mongo.Collection {
	return cr.cli.Database("courtDB").Collection("suspensions")
}
