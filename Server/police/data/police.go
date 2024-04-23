package data

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrafficPolice struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}

func (tp *TrafficPolice) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tp)
}

func (tp *TrafficPolice) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tp)
}
