package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mup struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name           string               `bson:"name" json:"name"`
	Address        Address              `bson:"address" json:"address"`
	Vehicles       []primitive.ObjectID `bson:"vehicles" json:"vehicles"`
	TrafficPermits []primitive.ObjectID `bson:"trafficPermits" json:"trafficPermits"`
	Plates         []string             `bson:"plates" json:"plates"`
	DrivingBans    []primitive.ObjectID `bson:"drivingBans" json:"drivingBans"`
	Registrations  []string             `bson:"registrations" json:"registrations"`
}

type DrivingBan struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Reason   string             `bson:"reason" json:"reason"`
	Duration time.Time          `bson:"duration" json:"duration"`
	Person   primitive.ObjectID `bson:"person" json:"person"`
}

type DrivingBans []DrivingBan

type TrafficPermit struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number         string             `bson:"number" json:"number"`
	IssuedDate     time.Time          `bson:"issuedDate" json:"issuedDate"`
	ExpirationDate time.Time          `bson:"expirationDate" json:"expirationDate"`
	Approved       bool               `bson:"approved" json:"approved"`
	Person         primitive.ObjectID `bson:"person" json:"person"`
}

type TrafficPermits []TrafficPermit

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

func (db *DrivingBans) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(db)
}

func (db *DrivingBans) FromJSON(r io.Reader) error {
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

func (tp *TrafficPermits) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tp)
}

func (tp *TrafficPermits) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tp)
}
