package data

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Court struct {
	ID                  primitive.ObjectID        `bson:"_id, omitempty" json:"id"`
	Name                string                    `bson:"name" json:"name"`
	Address             Address                   `bson:"address" json:"address"`
	HearingsPerson      []CourtHearingPerson      `bson:"hearingsPerson" json:"hearingsPerson"`
	HearingsLegalEntity []CourtHearingLegalEntity `bson:"hearingsLegalEntity" json:"hearingsLegalEntity"`
}

type CourtHearingPerson struct {
	ID       primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Reason   string             `bson:"reason" json:"reason"`
	DateTime time.Time          `bson:"dateTime" json:"dateTime"`
	Court    Court              `bson:"court" json:"court"`
	Person   Person             `bson:"person" json:"person"`
}

type CourtHearingLegalEntity struct {
	ID          primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Reason      string             `bson:"reason" json:"reason"`
	DateTime    time.Time          `bson:"dateTime" json:"dateTime"`
	Court       Court              `bson:"court" json:"court"`
	LegalEntity LegalEntity        `bson:"legalEntity" json:"legalEntity"`
}

func (c *Court) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Court) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}

func (chp *CourtHearingPerson) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(chp)
}

func (chp *CourtHearingPerson) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(chp)
}

func (chle *CourtHearingLegalEntity) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(chle)
}

func (chle *CourtHearingLegalEntity) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(chle)
}
