package data

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TIRE
const (
	Winter string = "WINTER"
	Summer string = "SUMMER"
)

type Vehicle struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Brand        string             `bson:"brand" json:"brand"`
	Model        string             `bson:"model" json:"model"`
	Year         int8               `bson:"year" json:"year"`
	Registration Registration       `bson:"registration" json:"registration"`
	Plates       Plates             `bson:"plates" json:"plates"`
	Owner        Person             `bson:"owner" json:"owner"`
	Tire         string             `bson:"tire" json:"tire"`
}

type Registration struct {
	RegistrationNumber string             `bson:"registration_number" json:"registration_number"`
	ExpirationDate     primitive.DateTime `bson:"expiration_date" json:"expiration_date"`
}

type Plates struct {
	RegistrationNumber string `bson:"registration_number" json:"registration_number"`
	PlateType          string `bson:"plate_type" json:"plate_type"`
}

func (p *Vehicle) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Vehicle) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *Registration) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Registration) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *Plates) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Plates) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
