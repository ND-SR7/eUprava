package data

import (
	"encoding/json"
	"io"
)

type TrafficViolation struct {
	Reason       string `bson:"reason" json:"reason"`
	Description  string `bson:"description" json:"description"`
	Time         string `bson:"time" json:"time"`
	Location     string `bson:"location" json:"location"`
	ViolatorJMBG string `bson:"violatorJMBG" json:"violatorJMBG"`
}

func (tv *TrafficViolation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tv)
}

func (tv *TrafficViolation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tv)
}
