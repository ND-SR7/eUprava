package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrafficPolice struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type TrafficViolation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ViolatorEmail string             `bson:"violatorEmail" json:"violatorEmail"`
	Reason        string             `bson:"reason" json:"reason"`
	Description   string             `bson:"description" json:"description"`
	Time          time.Time          `bson:"time" json:"time"`
	Location      string             `bson:"location" json:"location"`
}

type AlcoholTest struct {
	UserEmail    string  `bson:"userEmail" json:"userEmail"`
	AlcoholLevel float64 `bson:"alcoholLevel" json:"alcoholLevel"`
	Location     string  `bson:"location" json:"location"`
}

func (tp *TrafficPolice) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tp)
}

func (tp *TrafficPolice) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tp)
}

func (tv *TrafficViolation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tv)
}

func (tv *TrafficViolation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tv)
}

func (at *AlcoholTest) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(at)
}

func (at *AlcoholTest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(at)
}
