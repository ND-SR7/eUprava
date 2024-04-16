package data

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
	Id                string
	Email             string
	Password          string
	ActivationCode    string
	PasswordResetCode string
	Role              string
}

type Address struct {
	Municipality string
	Locality     string
	StreetName   string
	StreetNumber int
}

type Person struct {
	FirstName     string
	LastName      string
	Sex           string
	Citizenship   string
	DOB           string
	JMBG          string
	Account       Account
	Address       Address
	CourtHearings []CourtHearingPerson
}

type LegalEntity struct {
	Name          string
	Citizenship   string
	PIB           string
	MB            string
	Account       Account
	Address       Address
	CourtHearings []CourtHearingLegalEntity
}
