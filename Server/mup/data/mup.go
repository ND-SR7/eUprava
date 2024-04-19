package data

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
)

type Mup struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Address        Address            `bson:"address" json:"address"`
	TrafficPermits []TrafficPermit    `bson:"traffic_permits" json:"traffic_permits"`
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
	ExpirationDate primitive.DateTime `bson:"expirationDate" json:"expirationDate"`
	Person         Person             `bson:"person" json:"person"`
}

func (p *Mup) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Mup) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *DrivingBan) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *DrivingBan) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *TrafficPermit) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *TrafficPermit) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
