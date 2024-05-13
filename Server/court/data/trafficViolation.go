package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrafficViolation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Reason        string             `bson:"reason" json:"reason"`
	Description   string             `bson:"description" json:"description"`
	Time          time.Time          `bson:"time" json:"time"`
	Location      string             `bson:"location" json:"location"`
	ViolatorEmail string             `bson:"violatorEmail" json:"violatorEmail"`
}

func (tv *TrafficViolation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tv)
}

func (tv *TrafficViolation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tv)
}
