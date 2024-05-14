package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DrivingBan struct {
	Reason   string             `bson:"reason" json:"reason"`
	Duration time.Time          `bson:"duration" json:"duration"`
	Person   primitive.ObjectID `bson:"person" json:"person"`
}

func (db *DrivingBan) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(db)
}

func (db *DrivingBan) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(db)
}
