package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Warrant struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	TrafficViolation primitive.ObjectID `bson:"trafficViolation" json:"trafficViolation"`
	IssuedOn         time.Time          `bson:"issuedOn" json:"issuedOn"`
	IssuedFor        string             `bson:"issuedFor" json:"issuedFor"`
}

type Warrants []Warrant

func (a *Warrant) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Warrant) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}

func (a *Warrants) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Warrants) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
