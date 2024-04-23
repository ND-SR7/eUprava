package data

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Role
const (
	User  string = "USER"
	Admin string = "ADMIN"
)

// Sex
const (
	Male   string = "MALE"
	Female string = "FEMALE"
)

type Account struct {
	ID                primitive.ObjectID `bson:"_id, omitempty" json:"id"`
	Email             string             `bson:"email" json:"email"`
	Password          string             `bson:"password" json:"password"`
	ActivationCode    string             `bson:"activationCode" json:"activationCode"`
	PasswordResetCode string             `bson:"passwordResetCode" json:"passwordResetCode"`
	Role              string             `bson:"role" json:"role"`
}

type Address struct {
	Municipality string `bson:"municipality" json:"municipality"`
	Locality     string `bson:"locality" json:"locality"`
	StreetName   string `bson:"streetName" json:"streetName"`
	StreetNumber int    `bson:"streetNumber" json:"streetNumber"`
}

type Person struct {
	FirstName   string  `bson:"firstName" json:"firstName"`
	LastName    string  `bson:"lastName" json:"lastName"`
	Sex         string  `bson:"sex" json:"sex"`
	Citizenship string  `bson:"citizenship" json:"citizenship"`
	DOB         string  `bson:"dob" json:"dob"`
	JMBG        string  `bson:"jmbg" json:"jmbg"`
	Account     Account `bson:"account" json:"account"`
	Address     Address `bson:"address" json:"address"`
}

type LegalEntity struct {
	Name        string  `bson:"name" json:"name"`
	Citizenship string  `bson:"citizenship" json:"citizenship"`
	PIB         string  `bson:"pib" json:"pib"`
	MB          string  `bson:"mb" json:"mb"`
	Account     Account `bson:"account" json:"account"`
	Address     Address `bson:"address" json:"address"`
}

func (a *Account) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Account) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

func (a *Address) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(a)
}

func (a *Address) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(a)
}

func (p *Person) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Person) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (le *LegalEntity) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(le)
}

func (le *LegalEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(le)
}
