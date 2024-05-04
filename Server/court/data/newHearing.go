package data

import (
	"encoding/json"
	"io"
)

type NewCourtHearingPerson struct {
	Reason   string `bson:"reason" json:"reason"`
	DateTime string `bson:"dateTime" json:"dateTime"`
	Court    string `bson:"court" json:"court"`
	Person   string `bson:"person" json:"person"`
}

type NewCourtHearingLegalEntity struct {
	Reason      string `bson:"reason" json:"reason"`
	DateTime    string `bson:"dateTime" json:"dateTime"`
	Court       string `bson:"court" json:"court"`
	LegalEntity string `bson:"legalEntity" json:"legalEntity"`
}

type RescheduleCourtHearing struct {
	HearingID string `bson:"hearingID" json:"hearingID"`
	DateTime  string `bson:"dateTime" json:"dateTime"`
}

func (nh *NewCourtHearingPerson) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(nh)
}

func (nh *NewCourtHearingPerson) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(nh)
}

func (nh *NewCourtHearingLegalEntity) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(nh)
}

func (nh *NewCourtHearingLegalEntity) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(nh)
}
