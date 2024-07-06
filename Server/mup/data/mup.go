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
	Person   string             `bson:"person" json:"person"`
}

type DrivingBans []DrivingBan

type TrafficPermit struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number         string             `bson:"number" json:"number"`
	IssuedDate     time.Time          `bson:"issuedDate" json:"issuedDate"`
	ExpirationDate time.Time          `bson:"expirationDate" json:"expirationDate"`
	Approved       bool               `bson:"approved" json:"approved"`
	Person         string             `bson:"person" json:"person"`
}

type TrafficPermits []TrafficPermit

type TrafficPermitDetails struct {
	ID             primitive.ObjectID `json:"id"`
	Number         string             `json:"number"`
	IssuedDate     time.Time          `json:"issuedDate"`
	ExpirationDate time.Time          `json:"expirationDate"`
	Approved       bool               `json:"approved"`
	Person         string             `json:"person"`
	FirstName      string             `json:"firstName"`
	LastName       string             `json:"lastName"`
}

type TrafficPermitDetailsList []TrafficPermitDetails

type DrivingPermitDetails struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Number         string             `bson:"number" json:"number"`
	IssuedDate     time.Time          `bson:"issuedDate" json:"issuedDate"`
	ExpirationDate time.Time          `bson:"expirationDate" json:"expirationDate"`
	Approved       bool               `bson:"approved" json:"approved"`
	Person         string             `bson:"person" json:"person"`
	FirstName      string             `bson:"firstName" json:"firstName"`
	LastName       string             `bson:"lastName" json:"lastName"`
}

type DrivingPermitDetailsList []DrivingPermitDetails

func (dpd *DrivingPermitDetails) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(dpd)
}

func (dpd *DrivingPermitDetails) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(dpd)
}

func (dpdl *DrivingPermitDetailsList) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(dpdl)
}

func (dpdl *DrivingPermitDetailsList) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(dpdl)
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
