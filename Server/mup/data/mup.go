package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mup struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Address        Address            `bson:"address" json:"address"`
	TrafficPermits []TrafficPermit    `bson:"trafficPermits" json:"trafficPermits"`
	Plates         []Plates           `bson:"plates" json:"plates"`
	DrivingBans    []DrivingBan       `bson:"drivingBans" json:"drivingBans"`
	Registrations  []Registration     `bson:"registrations" json:"registrations"`
}

type DrivingBan struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Reason   string             `bson:"reason" json:"reason"`
	Duration string             `bson:"duration" json:"duration"`
	Person   Person             `bson:"person" json:"person"`
}

type TrafficPermit struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number         string             `bson:"number" json:"number"`
	ExpirationDate time.Time          `bson:"expirationDate" json:"expirationDate"`
	Person         Person             `bson:"person" json:"person"`
}

func (m *Mup) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(m)
}

func (m *Mup) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(m)
}

func (db *DrivingBan) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(db)
}

func (db *DrivingBan) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(db)
}

func (tp *TrafficPermit) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tp)
}

func (tp *TrafficPermit) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tp)
}
