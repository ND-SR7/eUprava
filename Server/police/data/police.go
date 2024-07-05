package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrafficPolice struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name" json:"name"`
}

type TrafficViolation struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ViolatorJMBG string             `bson:"violatorJMBG" json:"violatorJMBG"`
	Reason       string             `bson:"reason" json:"reason"`
	Description  string             `bson:"description" json:"description"`
	Time         time.Time          `bson:"time" json:"time"`
	Location     string             `bson:"location" json:"location"`
}

type DriverCheck struct {
	JMBG         string  `bson:"jmbg" json:"jmbg"`
	AlcoholLevel float64 `bson:"alcoholLevel" json:"alcoholLevel"`
	Tire         string  `bson:"tire" json:"tire"`
	PlatesNumber string  `bson:"platesNumber" json:"platesNumber"`
	Location     string  `bson:"location" json:"location"`
}

type AlcoholRequest struct {
	AlcoholLevel float64 `json:"alcoholLevel"`
	JMBG         string  `json:"jmbg"`
	Location     string  `json:"location"`
}

type DriverBanAndPermitRequest struct {
	JMBG     string `json:"jmbg"`
	Location string `json:"location"`
}

type VehicleTireCheck struct {
	TireType string `json:"tireType"`
	JMBG     string `json:"jmbg"`
	Location string `json:"location"`
}

type CheckVehicleRegistration struct {
	PlatesNumber string `json:"platesNumber"`
	JMBG         string `json:"jmbg"`
	Location     string `json:"location"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type DrivingBan struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Reason   string             `bson:"reason" json:"reason"`
	Duration time.Time          `bson:"duration" json:"duration"`
	Person   string             `bson:"person" json:"person"`
}

func (tp *TrafficPolice) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tp)
}

func (tp *TrafficPolice) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tp)
}

func (tv *TrafficViolation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tv)
}

func (tv *TrafficViolation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tv)
}

func (at *DriverCheck) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(at)
}

func (at *DriverCheck) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(at)
}

func (r *Response) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func (r *Response) FromJSON(re io.Reader) error {
	d := json.NewDecoder(re)
	return d.Decode(r)
}

func (r *DrivingBan) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(r)
}

func (r *DrivingBan) FromJSON(re io.Reader) error {
	d := json.NewDecoder(re)
	return d.Decode(r)
}
