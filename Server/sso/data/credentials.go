package data

import (
	"encoding/json"
	"io"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePassword struct {
	Email           string `json:"email"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func (c *Credentials) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Credentials) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(c)
}

func (cp *ChangePassword) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(cp)
}

func (cp *ChangePassword) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(cp)
}
