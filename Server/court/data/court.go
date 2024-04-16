package data

type Court struct {
	Id                  string
	Name                string
	Address             Address
	HearingsPerson      []CourtHearingPerson
	HearingsLegalEntity []CourtHearingLegalEntity
}

type CourtHearingPerson struct {
	Id       string
	Reason   string
	DateTime string
	Court    Court
	Person   Person
}

type CourtHearingLegalEntity struct {
	Id          string
	Reason      string
	DateTime    string
	Court       Court
	LegalEntity LegalEntity
}
