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

type CourtHearing interface {
	GetID() primitive.ObjectID
	GetReason() string
	GetDateTime() time.Time
	GetCourt() primitive.ObjectID
	GetSubjet() primitive.ObjectID
}

type CourtHearingPerson struct {
	ID       primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Reason   string             `bson:"reason" json:"reason"`
	DateTime time.Time          `bson:"dateTime" json:"dateTime"`
	Court    primitive.ObjectID `bson:"court" json:"court"`
	Person   primitive.ObjectID `bson:"person" json:"person"`
}

type CourtHearingLegalEntity struct {
	ID          primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Reason      string             `bson:"reason" json:"reason"`
	DateTime    time.Time          `bson:"dateTime" json:"dateTime"`
	Court       primitive.ObjectID `bson:"court" json:"court"`
	LegalEntity primitive.ObjectID `bson:"legalEntity" json:"legalEntity"`
}

func (c *Court) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Court) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}

func (chp *CourtHearingPerson) GetID() primitive.ObjectID {
	return chp.ID
}

func (chp *CourtHearingPerson) GetReason() string {
	return chp.Reason
}

func (chp *CourtHearingPerson) GetDateTime() time.Time {
	return chp.DateTime
}

func (chp *CourtHearingPerson) GetCourt() primitive.ObjectID {
	return chp.Court
}

func (chp *CourtHearingPerson) GetSubjet() primitive.ObjectID {
	return chp.Person
}

func (chp *CourtHearingPerson) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(chp)
}

func (chp *CourtHearingPerson) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(chp)
}

func (chle *CourtHearingLegalEntity) GetID() primitive.ObjectID {
	return chle.ID
}

func (chle *CourtHearingLegalEntity) GetReason() string {
	return chle.Reason
}

func (chle *CourtHearingLegalEntity) GetDateTime() time.Time {
	return chle.DateTime
}

func (chle *CourtHearingLegalEntity) GetCourt() primitive.ObjectID {
	return chle.Court
}

func (chle *CourtHearingLegalEntity) GetSubjet() primitive.ObjectID {
	return chle.LegalEntity
}

func (chle *CourtHearingLegalEntity) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(chle)
}

func (chle *CourtHearingLegalEntity) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(chle)
}
