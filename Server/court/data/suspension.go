package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Suspension struct {
	ID     primitive.ObjectID `bson:"_id" json:"id"`
	From   time.Time          `bson:"from" json:"from"`
	To     time.Time          `bson:"to" json:"to"`
	Person primitive.ObjectID `bson:"person" json:"person"`
}

type NewSuspension struct {
	From   string `bson:"from" json:"from"`
	To     string `bson:"to" json:"to"`
	Person string `bson:"person" json:"person"`
}

func (s *Suspension) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(s)
}

func (s *Suspension) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(s)
}

func (ns *NewSuspension) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ns)
}

func (ns *NewSuspension) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ns)
}
