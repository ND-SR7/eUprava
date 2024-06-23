package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstituteForStatistics struct {
	ID          primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Address     Address            `bson:"address" json:"address"`
	TrafficData []TrafficData      `bson:"trafficData" json:"trafficData"`
}

type StatisticsData struct {
	ID     primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Date   time.Time          `bson:"date" json:"date"`
	Region string             `bson:"region" json:"region"`
	Year   int                `bson:"year" json:"year"`
	Month  int                `bson:"month" json:"month"`
}

type TrafficData struct {
	ID primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	StatisticsData
	ViolationType string    `bson:"violationType" json:"violationType"`
	Vehicles      []Vehicle `bson:"vehicles" json:"vehicles"`
}

type Vehicle struct {
	ID           primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Brand        string             `bson:"brand" json:"brand"`
	Model        string             `bson:"model" json:"model"`
	Year         int                `bson:"year" json:"year"`
	Registration string             `bson:"registration" json:"registration"`
	Plates       string             `bson:"plates" json:"plates"`
	Owner        string             `bson:"owner" json:"owner"`
}

type Vehicles []Vehicle

type Registration struct {
	RegistrationNumber string             `bson:"registrationNumber" json:"registrationNumber"`
	IssuedDate         time.Time          `bson:"issuedDate" json:"issuedDate"`
	ExpirationDate     time.Time          `bson:"expirationDate" json:"expirationDate"`
	VehicleID          primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
	Approved           bool               `bson:"approved" json:"approved"`
}

type Plates struct {
	RegistrationNumber string             `bson:"registrationNumber" json:"registrationNumber"`
	PlatesNumber       string             `bson:"platesNumber" json:"platesNumber"`
	PlateType          string             `bson:"plateType" json:"plateType"`
	VehicleID          primitive.ObjectID `bson:"vehicleID" json:"vehicleID"`
}

type ListOfPlates []Plates

type TrafficViolation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ViolatorJMBG string             `bson:"violatorJMBG" json:"violatorJMBG"`
	Reason       string             `bson:"reason" json:"reason"`
	Description  string             `bson:"description" json:"description"`
	Time         time.Time          `bson:"time" json:"time"`
	Location     string             `bson:"location" json:"location"`
}

type TrafficViolations []TrafficViolation

func (i *InstituteForStatistics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

func (i *InstituteForStatistics) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(i)
}

func (sd *StatisticsData) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(sd)
}

func (sd *StatisticsData) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(sd)
}

func (td *TrafficData) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(td)
}

func (td *TrafficData) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(td)
}

func (v *Vehicle) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}

func (v *Vehicle) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)
}

func (v *Vehicles) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}

func (v *Vehicles) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)
}

func (re *Registration) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(re)
}

func (re *Registration) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(re)
}

func (p *Plates) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Plates) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (p *ListOfPlates) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *ListOfPlates) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (tv *TrafficViolation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tv)
}

func (tv *TrafficViolation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tv)
}

func (tv *TrafficViolations) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tv)
}

func (tv *TrafficViolations) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tv)
}
