package data

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type SSORepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// Constructor
func New(ctx context.Context, logger *log.Logger) (*SSORepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	sr := &SSORepo{
		cli:    client,
		logger: logger,
	}

	return sr, nil
}

// Disconnect
func (sr *SSORepo) Disconnect(ctx context.Context) error {
	err := sr.cli.Disconnect(ctx)
	if err != nil {
		sr.logger.Fatal(err.Error())
		return err
	}
	return nil
}

// Check database connection
func (sr *SSORepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping DB
	err := sr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		sr.logger.Println(err.Error())
	}

	// Print DBs
	databases, err := sr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		sr.logger.Println(err.Error())
	}
	sr.logger.Println(databases)
}

// Returns Account for specified email.
func (sr *SSORepo) FindAccountByEmail(email string) (Account, error) {
	persons := sr.getPersonsCollection()
	legalEntities := sr.getLegalEntitiesCollection()
	filter := bson.M{"account.email": email}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var person Person
	err := persons.FindOne(ctx, filter).Decode(&person)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			// Error other than not finding the account
			return Account{}, err
		}
		// Account not found in persons collection, try legal entities collection
		var legalEntity LegalEntity
		err = legalEntities.FindOne(ctx, filter).Decode(&legalEntity)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				// Account not found in both collections
				return Account{}, errors.New("account not found")
			}
			// Other error occurred
			return Account{}, err
		}

		return legalEntity.Account, nil
	}

	return person.Account, nil
}

// Getters for collections

func (sr *SSORepo) getPersonsCollection() *mongo.Collection {
	return sr.cli.Database("sso_db").Collection("persons")
}

func (sr *SSORepo) getLegalEntitiesCollection() *mongo.Collection {
	return sr.cli.Database("sso_db").Collection("legalEntities")
}
