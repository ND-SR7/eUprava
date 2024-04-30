package data

import (
	"encoding/json"
	"io"
)

type NewPerson struct {
	Email        string `bson:"email" json:"email"`
	Password     string `bson:"password" json:"password"`
	FirstName    string `bson:"firstName" json:"firstName"`
	LastName     string `bson:"lastName" json:"lastName"`
	Sex          string `bson:"sex" json:"sex"`
	Citizenship  string `bson:"citizenship" json:"citizenship"`
	DOB          string `bson:"dob" json:"dob"`
	JMBG         string `bson:"jmbg" json:"jmbg"`
	Role         string `bson:"role" json:"role"`
	Municipality string `bson:"municipality" json:"municipality"`
	Locality     string `bson:"locality" json:"locality"`
	StreetName   string `bson:"streetName" json:"streetName"`
	StreetNumber int    `bson:"streetNumber" json:"streetNumber"`
}

type NewLegalEntity struct {
	Email        string `bson:"email" json:"email"`
	Password     string `bson:"password" json:"password"`
	Name         string `bson:"name" json:"name"`
	Citizenship  string `bson:"citizenship" json:"citizenship"`
	PIB          string `bson:"pib" json:"pib"`
	MB           string `bson:"mb" json:"mb"`
	Role         string `bson:"role" json:"role"`
	Municipality string `bson:"municipality" json:"municipality"`
	Locality     string `bson:"locality" json:"locality"`
	StreetName   string `bson:"streetName" json:"streetName"`
	StreetNumber int    `bson:"streetNumber" json:"streetNumber"`
}

func (np *NewPerson) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(np)
}

func (np *NewPerson) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(np)
}

func (nle *NewLegalEntity) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(nle)
}

func (nle *NewLegalEntity) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(nle)
}
