package data

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Initialize database
func (sr *SSORepo) Initialize(ctx context.Context) error {
	db := sr.cli.Database("ssoDB")

	err := db.Collection("persons").Drop(ctx)
	if err != nil {
		return err
	}
	err = db.Collection("legalEntities").Drop(ctx)
	if err != nil {
		return err
	}

	persons := []Person{
		{
			FirstName:   "Mika",
			LastName:    "Mikic",
			DOB:         "2000-05-12",
			Sex:         "MALE",
			Citizenship: "serbian",
			JMBG:        "123456789",
			Account: Account{
				ID:                primitive.NewObjectID(),
				Email:             "mika@mail.com",
				Password:          "$2a$10$Q9ZCNMO5uGBAkQdMOhLLBe8JzvjK/oWjyq6.Tv8/O4UCpcW8.ymKS",
				ActivationCode:    uuid.New().String(),
				PasswordResetCode: uuid.New().String(),
				Role:              "ADMIN",
				Activated:         true,
			},
			Address: Address{
				Municipality: "Novi Sad",
				Locality:     "Novi Sad",
				StreetName:   "Somborska",
				StreetNumber: 9,
			},
		},
		{
			FirstName:   "Ana",
			LastName:    "Anic",
			DOB:         "2002-10-03",
			Sex:         "FEMALE",
			Citizenship: "serbian",
			JMBG:        "987654321",
			Account: Account{
				ID:                primitive.NewObjectID(),
				Email:             "ana@mail.com",
				Password:          "$2a$10$Lkjw/s5D9s1K38tSvsPTlOYI43ZJeu1c4.nMzH9nrod5Z1eKBRn4C",
				ActivationCode:    uuid.New().String(),
				PasswordResetCode: uuid.New().String(),
				Role:              "ADMIN",
				Activated:         true,
			},
			Address: Address{
				Municipality: "Novi Sad",
				Locality:     "Novi Sad",
				StreetName:   "Kamenicka",
				StreetNumber: 11,
			},
		},
	}

	var bsonPersons []interface{}
	for _, p := range persons {
		bsonPersons = append(bsonPersons, p)
	}

	_, err = db.Collection("persons").InsertMany(ctx, bsonPersons)
	if err != nil {
		return err
	}

	legalEntities := []LegalEntity{
		{
			Name:        "Test LE",
			Citizenship: "serbian",
			PIB:         "789123456",
			MB:          "00045698",
			Account: Account{
				ID:                primitive.NewObjectID(),
				Email:             "testle@mail.com",
				Password:          "$2a$10$Ddb59VQRIzEPyKHz1s1lc./0PyH1f5Z4Rz0psFwz75G80fxaBSsZS",
				ActivationCode:    "activate123",
				PasswordResetCode: "reset123",
				Role:              "ADMIN",
				Activated:         true,
			},
			Address: Address{
				Municipality: "Novi Sad",
				Locality:     "Novi Sad",
				StreetName:   "Jovan Zmaj",
				StreetNumber: 13,
			},
		},
		{
			Name:        "Legal Entity Inc.",
			Citizenship: "croatian",
			PIB:         "147369258",
			MB:          "88569536",
			Account: Account{
				ID:                primitive.NewObjectID(),
				Email:             "lei@mail.com",
				Password:          "$2a$10$Ddb59VQRIzEPyKHz1s1lc./0PyH1f5Z4Rz0psFwz75G80fxaBSsZS",
				ActivationCode:    "activateLEI",
				PasswordResetCode: "resetLEI",
				Role:              "ADMIN",
				Activated:         true,
			},
			Address: Address{
				Municipality: "Zagreb",
				Locality:     "Zagreb",
				StreetName:   "Kralja Petra",
				StreetNumber: 50,
			},
		},
	}

	var bsonLegalEntities []interface{}
	for _, le := range legalEntities {
		bsonLegalEntities = append(bsonLegalEntities, le)
	}

	_, err = db.Collection("legalEntities").InsertMany(ctx, bsonLegalEntities)
	if err != nil {
		return err
	}

	return nil
}

// Find person based on provided id
func (sr *SSORepo) GetPersonByID(id string) (Person, error) {
	persons := sr.getPersonsCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Person{}, err
	}

	filter := bson.M{"account._id": objID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var person Person
	err = persons.FindOne(ctx, filter).Decode(&person)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Person{}, errors.New("person not found")
	} else if err != nil {
		return Person{}, err
	}

	return person, nil
}

// Find legal entity based on provided id
func (sr *SSORepo) GetLegalEntityByID(id string) (LegalEntity, error) {
	legalEntities := sr.getLegalEntitiesCollection()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return LegalEntity{}, err
	}

	filter := bson.M{"account._id": objID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var legalEntity LegalEntity
	err = legalEntities.FindOne(ctx, filter).Decode(&legalEntity)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return LegalEntity{}, errors.New("person not found")
	} else if err != nil {
		return LegalEntity{}, err
	}

	return legalEntity, nil
}

// Find person based on provided email
func (sr *SSORepo) GetPersonByEmail(email string) (Person, error) {
	persons := sr.getPersonsCollection()

	filter := bson.M{"account.email": email}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var person Person
	err := persons.FindOne(ctx, filter).Decode(&person)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Person{}, errors.New("person not found")
	} else if err != nil {
		return Person{}, err
	}

	// Removing sensitive data
	person.Account.Password = ""
	person.Account.ActivationCode = ""
	person.Account.PasswordResetCode = ""

	return person, nil
}

// Find legal entity based on provided email
func (sr *SSORepo) GetLegalEntityByEmail(email string) (LegalEntity, error) {
	legalEntities := sr.getLegalEntitiesCollection()

	filter := bson.M{"account.email": email}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var legalEntity LegalEntity
	err := legalEntities.FindOne(ctx, filter).Decode(&legalEntity)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return LegalEntity{}, errors.New("legal entity not found")
	} else if err != nil {
		return LegalEntity{}, err
	}

	// Removing sensitive data
	legalEntity.Account.Password = ""
	legalEntity.Account.ActivationCode = ""
	legalEntity.Account.PasswordResetCode = ""

	return legalEntity, nil
}

// Find person based on provided JMBG
func (sr *SSORepo) GetPersonByJMBG(jmbg string) (Person, error) {
	persons := sr.getPersonsCollection()

	filter := bson.M{"jmbg": jmbg}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var person Person
	err := persons.FindOne(ctx, filter).Decode(&person)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return Person{}, errors.New("person not found")
	} else if err != nil {
		return Person{}, err
	}

	// Removing sensitive data
	person.Account.Password = ""
	person.Account.ActivationCode = ""
	person.Account.PasswordResetCode = ""

	return person, nil
}

// Find legal entity based on provided MB
func (sr *SSORepo) GetLegalEntityByMB(mb string) (LegalEntity, error) {
	legalEntities := sr.getLegalEntitiesCollection()

	filter := bson.M{"mb": mb}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var legalEntity LegalEntity
	err := legalEntities.FindOne(ctx, filter).Decode(&legalEntity)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return LegalEntity{}, errors.New("legal entity not found")
	} else if err != nil {
		return LegalEntity{}, err
	}

	// Removing sensitive data
	legalEntity.Account.Password = ""
	legalEntity.Account.ActivationCode = ""
	legalEntity.Account.PasswordResetCode = ""

	return legalEntity, nil
}

// Inserts new person into collection
func (sr *SSORepo) CreatePerson(newPerson NewPerson) error {
	emailOK := true
	_, err := sr.FindAccountByEmail(newPerson.Email)
	if err == nil {
		emailOK = false
	}

	passwordOK := CheckPassword(newPerson.Password)

	if emailOK && passwordOK {
		hashedPassword, err := HashPassword(newPerson.Password)
		if err != nil {
			log.Fatalf("Error while hashing password: %s", err.Error())
			return err
		}

		collection := sr.getPersonsCollection()

		person := Person{
			FirstName:   newPerson.FirstName,
			LastName:    newPerson.LastName,
			Sex:         newPerson.Sex,
			Citizenship: newPerson.Citizenship,
			DOB:         newPerson.DOB,
			JMBG:        newPerson.JMBG,
			Account: Account{
				ID:                primitive.NewObjectID(),
				Email:             newPerson.Email,
				Password:          hashedPassword,
				ActivationCode:    uuid.New().String(),
				PasswordResetCode: uuid.New().String(),
				Role:              newPerson.Role,
				Activated:         false,
			},
			Address: Address{
				Municipality: newPerson.Municipality,
				Locality:     newPerson.Locality,
				StreetName:   newPerson.StreetName,
				StreetNumber: newPerson.StreetNumber,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = collection.InsertOne(ctx, person)
		if err != nil {
			log.Fatalf("Failed to insert new person: %s", err.Error())
			return err
		}

		emailSent := SendEmail(person.Account.Email, person.Account.ActivationCode, "ACTIVATION")
		if !emailSent {
			log.Fatalln("Failed to send activation email")
			return errors.New("Failed to send email")
		}
		sr.logger.Println("Activation email sent")
	} else if !emailOK {
		return errors.New("email already in use")
	} else if !passwordOK {
		return errors.New("choose a stronger password")
	}

	return nil
}

// Inserts new legal entity into collection
func (sr *SSORepo) CreateLegalEntity(newLegalEntity NewLegalEntity) error {
	emailOK := true
	_, err := sr.FindAccountByEmail(newLegalEntity.Email)
	if err == nil {
		emailOK = false
	}

	passwordOK := CheckPassword(newLegalEntity.Password)

	if emailOK && passwordOK {
		hashedPassword, err := HashPassword(newLegalEntity.Password)
		if err != nil {
			log.Fatalf("Error while hashing password: %s", err.Error())
			return err
		}

		collection := sr.getLegalEntitiesCollection()

		legalEntity := LegalEntity{
			Name:        newLegalEntity.Name,
			Citizenship: newLegalEntity.Citizenship,
			PIB:         newLegalEntity.PIB,
			MB:          newLegalEntity.MB,
			Account: Account{
				ID:                primitive.NewObjectID(),
				Email:             newLegalEntity.Email,
				Password:          hashedPassword,
				ActivationCode:    uuid.New().String(),
				PasswordResetCode: uuid.New().String(),
				Role:              newLegalEntity.Role,
				Activated:         false,
			},
			Address: Address{
				Municipality: newLegalEntity.Municipality,
				Locality:     newLegalEntity.Locality,
				StreetName:   newLegalEntity.StreetName,
				StreetNumber: newLegalEntity.StreetNumber,
			},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = collection.InsertOne(ctx, legalEntity)
		if err != nil {
			log.Fatalf("Failed to insert new legal entity: %s", err.Error())
			return err
		}

		emailSent := SendEmail(legalEntity.Account.Email, legalEntity.Account.ActivationCode, "ACTIVATION")
		if !emailSent {
			log.Fatalln("Failed to send activation email")
			return errors.New("Failed to send email")
		}
		sr.logger.Println("Activation email sent")
	} else if !emailOK {
		return errors.New("email already in use")
	} else if !passwordOK {
		return errors.New("choose a stronger password")
	}

	return nil
}

// Updates activated flag in account to true for specified activation code
func (sr *SSORepo) ActivateAccount(activationCode string) error {
	persons := sr.getPersonsCollection()
	legalEntities := sr.getLegalEntitiesCollection()
	filter := bson.M{"account.activationCode": activationCode}

	update := bson.M{
		"$set": bson.M{
			"account.activated": true,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := persons.UpdateOne(ctx, filter, update)
	if err != nil {
		sr.logger.Println("Failed to activate account")
		return err
	}
	if result.ModifiedCount > 0 {
		sr.logger.Println("Successfully activated person's account")
		return nil
	}

	result, err = legalEntities.UpdateOne(ctx, filter, update)
	if err != nil {
		sr.logger.Println("Failed to activate account")
		return err
	}
	if result.ModifiedCount == 0 {
		sr.logger.Printf("Invalid activation code: %s", activationCode)
		return errors.New("invalid activation code: " + activationCode)
	}

	return nil
}

// Resets password for specified reset code
func (sr *SSORepo) ResetPassword(passwordResetCode, newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	persons := sr.getPersonsCollection()
	legalEntities := sr.getLegalEntitiesCollection()
	filter := bson.M{"account.passwordResetCode": passwordResetCode}
	update := bson.M{
		"$set": bson.M{
			"account.password":          hashedPassword,
			"account.passwordResetCode": uuid.New().String(),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := persons.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	} else if result.ModifiedCount == 0 {
		result, err = legalEntities.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		} else if result.MatchedCount == 0 {
			return errors.New("no account with provided reset code")
		}
	}

	return nil
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
	return sr.cli.Database("ssoDB").Collection("persons")
}

func (sr *SSORepo) getLegalEntitiesCollection() *mongo.Collection {
	return sr.cli.Database("ssoDB").Collection("legalEntities")
}
