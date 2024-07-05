package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tire type
const (
	Winter string = "WINTER"
	Summer string = "SUMMER"
)

type Vehicle struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Brand        string             `bson:"brand" json:"brand"`
	Model        string             `bson:"model" json:"model"`
	Year         int                `bson:"year" json:"year"`
	Registration Registration       `bson:"registration" json:"registration"`
	Plates       Plates             `bson:"plates" json:"plates"`
	Owner        Person             `bson:"owner" json:"owner"`
	Tire         string             `bson:"tire" json:"tire"`
}

type Vehicles []Vehicles

type Registration struct {
	RegistrationNumber string             `bson:"registrationNumber" json:"registrationNumber"`
	IssuedDate         time.Time          `bson:"issuedDate" json:"issuedDate"`
	ExpirationDate     time.Time          `bson:"expirationDate" json:"expirationDate"`
	VehicleID          primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
	Owner              string             `bson:"owner" json:"owner"`
	Plates             string             `bson:"plates" json:"plates"`
	Approved           bool               `bson:"approved" json:"approved"`
}

type Plates struct {
	RegistrationNumber string             `bson:"registrationNumber" json:"registrationNumber"`
	PlatesNumber       string             `bson:"platesNumber" json:"platesNumber"`
	PlateType          string             `bson:"plateType" json:"plateType"`
	Owner              string             `bson:"owner" json:"owner"`
	VehicleID          primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
}

type PlateRequest struct {
	Plate string `json:"plates"`
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
