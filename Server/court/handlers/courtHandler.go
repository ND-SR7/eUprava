package handlers

import (
	"court/data"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type CourtHandler struct {
	repo *data.CourtRepo
}

const InvalidRequestBody = "Invalid request body"
const InvalidRequestBodyError = "Error while decoding body"

// Constructor
func NewCourtHandler(r *data.CourtRepo) *CourtHandler {
	return &CourtHandler{r}
}

// Handler methods

// Retrieves hearing based on provided ID
func (ch *CourtHandler) GetCourtHearingByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hearingID := params["id"]

	log.Printf("Retrieving hearing with id '%s'", hearingID)

	hearing, err := ch.getHearing(hearingID)
	if err != nil {
		http.Error(w, "Failed to retrieve court hearing", http.StatusInternalServerError)
		log.Printf("Failed to retrieve court hearing: %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(hearing); err != nil {
		http.Error(w, "Error while encoding body", http.StatusInternalServerError)
		log.Printf("Error while encoding court hearing: %s", err.Error())
	}

	log.Println("Successfully retrieved requested hearing")
}

// Creates a new hearing for a person
func (ch *CourtHandler) CreateHearingPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating a new court hearing for person")

	var newHearing data.NewCourtHearingPerson
	if err := json.NewDecoder(r.Body).Decode(&newHearing); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println(InvalidRequestBodyError)
		return
	}

	err := ch.repo.CreateHearingPerson(newHearing)
	if err != nil {
		http.Error(w, "Failed to create new court hearing", http.StatusInternalServerError)
		log.Printf("Failed to create new court hearing: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Successfully created a new court hearing")
}

// Creates a new hearing for a legal entity
func (ch *CourtHandler) CreateHearingLegalEntity(w http.ResponseWriter, r *http.Request) {
	log.Println("Creating a new court hearing for legal entity")

	var newHearing data.NewCourtHearingLegalEntity
	if err := json.NewDecoder(r.Body).Decode(&newHearing); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println(InvalidRequestBodyError)
		return
	}

	err := ch.repo.CreateHearingLegalEntity(newHearing)
	if err != nil {
		http.Error(w, "Failed to create new court hearing", http.StatusInternalServerError)
		log.Printf("Failed to create new court hearing: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	log.Println("Successfully created a new court hearing")
}

// Reschedules already existing court hearing for a person
func (ch *CourtHandler) UpdateHearingPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("Rescheduling court hearing for a person")

	var rescheduledHearing data.RescheduleCourtHearing
	if err := json.NewDecoder(r.Body).Decode(&rescheduledHearing); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println(InvalidRequestBodyError)
		return
	}

	rescheduledDateTime, err := time.Parse("2006-01-02T15:04:05", rescheduledHearing.DateTime)
	if err != nil {
		http.Error(w, "Error while decoding rescheduled date and time", http.StatusBadRequest)
		log.Println(InvalidRequestBodyError)
		return
	}

	courtHearing, err := ch.getHearing(rescheduledHearing.HearingID)
	if err != nil {
		http.Error(w, "Failed to retrieve court hearing", http.StatusInternalServerError)
		log.Printf("Failed to retrieve court hearing: %s", err.Error())
		return
	}

	if rescheduledDateTime.Before(courtHearing.GetDateTime()) {
		http.Error(w, "Court hearing can't be rescheduled before set date and time", http.StatusBadRequest)
		log.Println(InvalidRequestBodyError)
		return
	}

	err = ch.repo.RescheduleCourtHearingPerson(rescheduledHearing)
	if err != nil {
		http.Error(w, "Failed to update court hearing", http.StatusInternalServerError)
		log.Printf("Failed to update court hearing: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Successfully rescheduled court hearing")
}

// Reschedules already existing court hearing for a legal entity
func (ch *CourtHandler) UpdateHearingLegalEntity(w http.ResponseWriter, r *http.Request) {
	log.Println("Rescheduling court hearing for a legal entity")

	var rescheduledHearing data.RescheduleCourtHearing
	if err := json.NewDecoder(r.Body).Decode(&rescheduledHearing); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println(InvalidRequestBodyError)
		return
	}

	rescheduledDateTime, err := time.Parse("2006-01-02T15:04:05", rescheduledHearing.DateTime)
	if err != nil {
		http.Error(w, "Error while decoding rescheduled date and time", http.StatusBadRequest)
		log.Println(InvalidRequestBodyError)
		return
	}

	courtHearing, err := ch.getHearing(rescheduledHearing.HearingID)
	if err != nil {
		http.Error(w, "Failed to retrieve court hearing", http.StatusInternalServerError)
		log.Printf("Failed to retrieve court hearing: %s", err.Error())
		return
	}

	if rescheduledDateTime.Before(courtHearing.GetDateTime()) {
		http.Error(w, "Court hearing can't be rescheduled before set date and time", http.StatusBadRequest)
		log.Println(InvalidRequestBodyError)
		return
	}

	err = ch.repo.RescheduleCourtHearingLegalEntity(rescheduledHearing)
	if err != nil {
		http.Error(w, "Failed to update court hearing", http.StatusInternalServerError)
		log.Printf("Failed to update court hearing: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Successfully rescheduled court hearing")
}

// Helper function for parsing court hearing interface into structs.
// Retrieves court hearing from repo and converts it
func (ch *CourtHandler) getHearing(id string) (data.CourtHearing, error) {
	hearing, err := ch.repo.GetHearingByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	if hearingPerson, ok := hearing.(*data.CourtHearingPerson); ok {
		return hearingPerson, nil
	} else if hearingLegalEntity, ok := hearing.(*data.CourtHearingLegalEntity); ok {
		return hearingLegalEntity, nil
	}

	return nil, fmt.Errorf("could not convert retrieved court hearing to any type")
}
