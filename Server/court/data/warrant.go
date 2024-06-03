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

type NewWarrant struct {
	TrafficViolation string `bson:"trafficViolation" json:"trafficViolation"`
	IssuedFor        string `bson:"issuedFor" json:"issuedFor"`
}

func (wa *Warrant) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(wa)
}

func (wa *Warrant) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(wa)
}

func (wa *Warrants) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(wa)
}

func (wa *Warrants) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(wa)
}

func (nw *NewWarrant) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(nw)
}

func (nw *NewWarrant) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(nw)
}
