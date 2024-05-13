package data

import (
	"encoding/json"
	"io"
)

type Warrant struct {
	ID               string `bson:"_id" json:"id"`
	TrafficViolation string `bson:"trafficViolation" json:"trafficViolation"`
	IssuedOn         string `bson:"issuedOn" json:"issuedOn"`
}

func (a *Warrant) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Warrant) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(a)
}
