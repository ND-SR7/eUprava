package handlers

import (
	"context"
	"court/clients"
	"court/data"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CourtHandler struct {
	repo *data.CourtRepo
	sso  clients.SSOClient
}

var secretKey = []byte("eUpravaT2")

const InvalidRequestBody = "Invalid request body"
const InvalidRequestBodyError = "Error while decoding body"

// Constructor
func NewCourtHandler(r *data.CourtRepo, s clients.SSOClient) *CourtHandler {
	return &CourtHandler{r, s}
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

func (ch *CourtHandler) RecieveCrimeReport(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved crime report")

	var trafficViolation data.TrafficViolation
	if err := json.NewDecoder(r.Body).Decode(&trafficViolation); err != nil {
		http.Error(w, InvalidRequestBody, http.StatusBadRequest)
		log.Println(InvalidRequestBody)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 4*time.Second)
	defer cancel()

	token := ch.extractTokenFromHeader(r)
	person, err := ch.sso.GetPersonByEmail(ctx, trafficViolation.ViolatorEmail, token)
	if err != nil {
		http.Error(w, "Error with services communication", http.StatusInternalServerError)
		log.Printf("Error while communicating with SSO service: %s", err.Error())
		return
	}

	courtHearing := data.NewCourtHearingPerson{
		Reason:   trafficViolation.Reason,
		DateTime: time.Now().Add(72 * time.Hour).String(),
		Court:    primitive.NewObjectID().Hex(), // TODO
		Person:   person.Account.ID.Hex(),
	}

	err = ch.repo.CreateHearingPerson(courtHearing)
	if err != nil {
		http.Error(w, "Error while creating new hearing", http.StatusInternalServerError)
		log.Printf("Error while creating new person hearing: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Successfully scheduled court hearing after crime report")
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

// JWT middleware
func (ch *CourtHandler) AuthorizeRoles(allowedRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
			tokenString := ch.extractTokenFromHeader(rr)
			if tokenString == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_, ok1 := claims["sub"].(string)
			role, ok2 := claims["role"].(string)
			if !ok1 || !ok2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			for _, allowedRole := range allowedRoles {
				if allowedRole == role {
					next.ServeHTTP(w, rr)
					return
				}
			}

			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

// Returns token string found in header, otherwise empty string
func (ch *CourtHandler) extractTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token != "" {
		return token[len("Bearer "):]
	}
	return ""
}
