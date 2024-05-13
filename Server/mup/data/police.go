package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"time"
)

type Warrant struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	TrafficViolation primitive.ObjectID `bson:"trafficViolation" json:"trafficViolation"`
	IssuedOn         time.Time          `bson:"issuedOn" json:"issuedOn"`
}

func (a *Warrant) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Warrant) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
