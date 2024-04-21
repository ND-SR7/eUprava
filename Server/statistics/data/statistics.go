package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InstituteForStatistics struct {
    Id         primitive.ObjectID 				`bson:"_id, omitempty" json:"id"`
    Name       string             				`bson:"name" json:"name"`
    Address    Address            				`bson:"address" json:"address"`
    CrimeData  []CrimeData        				`bson:"crimeData" json:"crimeData"`
    TrafficData []TrafficData     				`bson:"trafficData" json:"trafficData"`
}

type StatisticsData struct {
    Id     primitive.ObjectID 				`bson:"_id, omitempty" json:"id"`
    Date   time.Time          				`bson:"date" json:"date"`
    Region string             				`bson:"region" json:"region"`
    Year   int                				`bson:"year" json:"year"`
    Month  int                				`bson:"month" json:"month"`
}

type CrimeData struct {
    StatisticsData
    CrimeType string 				`bson:"crimeType" json:"crimeType"`
}

type TrafficData struct {
    StatisticsData
    ViolationType string    				`bson:"violationType" json:"violationType"`
    Vehicles      []Vehicle 				`bson:"vehicles" json:"vehicles"`
}

type Vehicle struct {
    Id             primitive.ObjectID 				`bson:"_id, omitempty" json:"id"`
    Brand          string             				`bson:"brand" json:"brand"`
    Model          string             				`bson:"model" json:"model"`
    Year           string             				`bson:"year" json:"year"`
    Registration   Registration       				`bson:"registration" json:"registration"`
}

type Registration struct {
    RegistrationNumber string    				`bson:"registrationNumber" json:"registrationNumber"`
    ExpiryDate         time.Time 				`bson:"expiryDate" json:"expiryDate"`
}

func (i *InstituteForStatistics) ToJSON(w io.Writer) error {
    e := json.NewEncoder(w)
    return e.Encode(i)
}

func (i *InstituteForStatistics) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(i)
}

func (s *StatisticsData) ToJSON(w io.Writer) error {
    e := json.NewEncoder(w)
    return e.Encode(s)
}

func (s *StatisticsData) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(s)
}

func (c *CrimeData) ToJSON(w io.Writer) error {
    e := json.NewEncoder(w)
    return e.Encode(c)
}

func (c *CrimeData) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(c)
}

func (t *TrafficData) ToJSON(w io.Writer) error {
    e := json.NewEncoder(w)
    return e.Encode(t)
}

func (t *TrafficData) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(t)
}

func (v *Vehicle) ToJSON(w io.Writer) error {
    e := json.NewEncoder(w)
    return e.Encode(v)
}

func (v *Vehicle) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(v)
}

func (re *Registration) ToJSON(w io.Writer) error {
    e := json.NewEncoder(w)
    return e.Encode(re)
}

func (re *Registration) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(re)
}
