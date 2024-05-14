package data

import (
	"encoding/json"
	"io"
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
	Owner        primitive.ObjectID `bson:"owner" json:"owner"`
}

type Vehicles []Vehicle

type Registration struct {
	RegistrationNumber string             `bson:"registrationNumber" json:"registrationNumber"`
	IssuedDate         time.Time          `bson:"issuedDate" json:"issuedDate"`
	ExpirationDate     time.Time          `bson:"expirationDate" json:"expirationDate"`
	VehicleID          primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
	Approved           bool               `bson:"approved" json:"approved"`
}

type Registrations []Registration

type Plates struct {
	RegistrationNumber string             `bson:"registrationNumber" json:"registrationNumber"`
	PlatesNumber       string             `bson:"platesNumber" json:"platesNumber"`
	PlateType          string             `bson:"plateType" json:"plateType"`
	VehicleID          primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
}

type ListOfPlates []Plates

func (v *Vehicle) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}

func (v *Vehicle) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)
}

func (v *Vehicles) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}

func (v *Vehicles) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)
}

func (re *Registration) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(re)
}

func (re *Registration) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(re)
}

func (re *Registrations) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(re)
}

func (re *Registrations) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(re)
}

func (p *Plates) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Plates) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *ListOfPlates) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *ListOfPlates) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
