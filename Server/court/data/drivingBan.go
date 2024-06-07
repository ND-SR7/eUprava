package data

import (
	"encoding/json"
	"io"
	"time"
)

type DrivingBan struct {
	Reason   string    `bson:"reason" json:"reason"`
	Duration time.Time `bson:"duration" json:"duration"`
	Person   string    `bson:"person" json:"person"`
}

func (db *DrivingBan) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(db)
}

func (db *DrivingBan) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(db)
}
