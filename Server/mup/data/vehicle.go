package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vehicle struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Brand        string             `bson:"brand" json:"brand"`
	Model        string             `bson:"model" json:"model"`
	Year         int                `bson:"year" json:"year"`
	Registration string             `bson:"registration" json:"registration"`
	Plates       string             `bson:"plates" json:"plates"`
	Owner        string             `bson:"owner" json:"owner"`
}

type Vehicles []Vehicle

type Registration struct {
	RegistrationNumber string             `bson:"registrationNumber" json:"registrationNumber"`
	IssuedDate         time.Time          `bson:"issuedDate" json:"issuedDate"`
	ExpirationDate     time.Time          `bson:"expirationDate" json:"expirationDate"`
	VehicleID          primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
	Owner              string             `bson:"owner" json:"owner"`
	Plates             string             `bson:"plates" json:"plates"`
	Approved           bool               `bson:"approved" json:"approved"`
}

type Registrations []Registration

type RegistrationDetails struct {
	RegistrationNumber string             `json:"registrationNumber"`
	IssuedDate         time.Time          `json:"issuedDate"`
	ExpirationDate     time.Time          `json:"expirationDate"`
	VehicleID          primitive.ObjectID `json:"vehicleID"`
	Owner              string             `json:"owner"`
	Plates             string             `json:"plates"`
	Approved           bool               `json:"approved"`
	FirstName          string             `json:"firstName"`
	LastName           string             `json:"lastName"`
	VehicleBrand       string             `json:"vehicleBrand"`
	VehicleModel       string             `json:"vehicleModel"`
}

type RegistrationDetailsList []RegistrationDetails

type Plates struct {
	RegistrationNumber string             `bson:"registrationNumber" json:"registrationNumber"`
	PlatesNumber       string             `bson:"platesNumber" json:"platesNumber"`
	PlateType          string             `bson:"plateType" json:"plateType"`
	Owner              string             `bson:"owner" json:"owner"`
	VehicleID          primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
}

type ListOfPlates []Plates

type PlateRequest struct {
	Plate string `json:"plates"`
}

type VehicleDTO struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Brand        string             `bson:"brand" json:"brand"`
	Model        string             `bson:"model" json:"model"`
	Year         int                `bson:"year" json:"year"`
	Registration Registration       `bson:"registration" json:"registration"`
	Plates       Plates             `bson:"plates" json:"plates"`
	Owner        string             `bson:"owner" json:"owner"`
}

type VehiclesDTO []VehicleDTO

// JSON methods...
