package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mup/clients"
	"mup/data"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ApplicationJson = "application/json"
const ContentType = "Content-Type"
const FailedToEncodePerson = "Failed to encode person"
const FailedToReadUsernameFromToken = "Failed to read username from token"
const FailedToEcnodeDrivingBans = "Failed to encode driving bans"
const InvalidID = "Invalid ID"
const FailedToDecodeRequestBody = "Failed to decode request body"

type KeyProduct struct{}

var secretKey = []byte("eUpravaT2")

type MupHandler struct {
	repo   *data.MUPRepo
	logger *log.Logger
	ssoc   clients.SSOClient
	cc     clients.CourtClient
}

func NewMupHandler(r *data.MUPRepo, log *log.Logger, ssoc clients.SSOClient,
	cc clients.CourtClient) *MupHandler {
	return &MupHandler{r, log, ssoc, cc}
}

// Ping
func (mh *MupHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Pong"))

	w.WriteHeader(http.StatusOK)
}

// TODO: implement check for persons driving licience
// TODO: implement check for persons registrations

//GET

func (mh *MupHandler) CheckForPersonsDrivingBans(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	personID, _ := primitive.ObjectIDFromHex("607d22b837ede6b71eef3e11")

	drivingBans, err := mh.repo.CheckForPersonsDrivingBans(ctx, personID)
	if err != nil {
		http.Error(rw, "Failed to retrieve persons driving bans", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(drivingBans); err != nil {
		http.Error(rw, FailedToEcnodeDrivingBans, http.StatusInternalServerError)
	}
	fmt.Println("Successfully fetched driving bans")
}

func (mh *MupHandler) CheckForRegisteredVehicles(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicles, err := mh.repo.RetrieveRegisteredVehicles(ctx)
	if err != nil {
		http.Error(rw, "Failed to retrieve registered vehicles", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(vehicles); err != nil {
		http.Error(rw, FailedToEcnodeDrivingBans, http.StatusInternalServerError)
	}
	fmt.Println("Successfully fetched registered vehicles")
}

// POST

func (mh *MupHandler) SubmitRegistrationRequest(rw http.ResponseWriter, r *http.Request) {
	var registration data.Registration

	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := mh.repo.SubmitRegistrationRequest(r.Context(), &registration); err != nil {
		log.Printf("Failed to submt registration request: %v", err)
		http.Error(rw, "Failed to submit registration request", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(registration); err != nil {
		log.Printf("Failed to encode registration request: %v", err)
		http.Error(rw, "Failed to encode registration request", http.StatusInternalServerError)
	}

	log.Printf("Successfully created registration request with id '%s'", registration.RegistrationNumber)
}

func (mh *MupHandler) SubmitTrafficPermitRequest(rw http.ResponseWriter, r *http.Request) {
	var trafficPermit data.TrafficPermit
	ctx := r.Context()
	tokenStr := mh.extractTokenFromHeader(r)

	jmbg, err := mh.getJMBGFromToken(tokenStr)
	if err != nil {
		fmt.Printf("Error while reading email from token: %v", err)
		http.Error(rw, FailedToReadUsernameFromToken, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&trafficPermit); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	user, err := mh.ssoc.GetUserByJMBG(ctx, jmbg, tokenStr)
	if err != nil {
		http.Error(rw, "Failed to get user by email from sso", http.StatusBadRequest)
		log.Printf("Failed to get user by email from sso: %v", err)
		return
	}

	warrants, err := mh.cc.CheckForPersonsWarrant(ctx, trafficPermit.Person, tokenStr)
	if err != nil {
		http.Error(rw, "Failed to get warrant from court service", http.StatusBadRequest)
		log.Printf("Failed to get warrant from court service: %v", err)
		return
	}

	if len(warrants) != 0 {
		http.Error(rw, "User is on warrant list", http.StatusBadRequest)
		log.Printf("User is on warrant list")
		return
	}

	trafficPermit.Person = user.JMBG
	trafficPermit.Approved = false
	trafficPermit.IssuedDate = time.Now()

	if err := mh.repo.SubmitTrafficPermitRequest(ctx, &trafficPermit); err != nil {
		log.Printf("Failed to submt traffic permit request: %v", err)
		http.Error(rw, "Failed to submit traffic permit request", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(trafficPermit); err != nil {
		log.Printf("Failed to encode registration request: %v", err)
		http.Error(rw, "Failed to encode traffic permit request", http.StatusInternalServerError)
	}

	log.Printf("Successfully created traffic permit request with id '%s'", trafficPermit.ID.Hex())
}

func (mh *MupHandler) SaveVehicle(rw http.ResponseWriter, r *http.Request) {
	var vehicle data.Vehicle

	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := mh.repo.SaveVehicle(r.Context(), &vehicle); err != nil {
		log.Printf("Failed to save vehicle: %v", err)
		http.Error(rw, "Failed to save vehicle", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(vehicle); err != nil {
		log.Printf("Failed to encode vehicle: %v", err)
		http.Error(rw, "Failed to encode vehicle", http.StatusInternalServerError)
	}

	log.Printf("Successfully created vehicle with id '%s'", vehicle.ID.Hex())
}

func (mh *MupHandler) IssueDrivingBan(rw http.ResponseWriter, r *http.Request) {
	var drivingBan data.DrivingBan

	if err := json.NewDecoder(r.Body).Decode(&drivingBan); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := mh.repo.IssueDrivingBan(r.Context(), &drivingBan); err != nil {
		log.Printf("Failed to save driving ban: %v", err)
		http.Error(rw, "Failed to driving ban", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(drivingBan); err != nil {
		log.Printf("Failed to encode driving ban: %v", err)
		http.Error(rw, "Failed to encode driving ban", http.StatusInternalServerError)
	}

	log.Printf("Successfully created driving ban with id '%s'", drivingBan.ID.Hex())
}

func (mh *MupHandler) ApproveRegistration(rw http.ResponseWriter, r *http.Request) {
	var registration data.Registration

	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := mh.repo.ApproveRegistration(r.Context(), registration); err != nil {
		log.Printf("Failed to approve registration: %v", err)
		http.Error(rw, "Failed to approve registration", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(registration); err != nil {
		log.Printf("Failed to encode approved registration: %v", err)
		http.Error(rw, "Failed to encode approved registration", http.StatusInternalServerError)
	}
	log.Printf("Successfully updated registration '%s'", registration.RegistrationNumber)
}

func (mh *MupHandler) ApproveTrafficPermitRequest(rw http.ResponseWriter, r *http.Request) {
	var trafficPermit data.TrafficPermit

	if err := json.NewDecoder(r.Body).Decode(&trafficPermit); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := mh.repo.ApproveTrafficPermitRequest(r.Context(), trafficPermit.ID); err != nil {
		log.Printf("Failed to approve traffic permit: %v", err)
		http.Error(rw, "Failed to approve traffic permit", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(trafficPermit); err != nil {
		log.Printf("Failed to encode traffic permit: %v", err)
		http.Error(rw, "Failed to encode approved traffic permit", http.StatusInternalServerError)
	}
	log.Printf("Successfully updated traffic permit '%s'", trafficPermit.ID.Hex())
}

// Save mup
func (mh *MupHandler) SaveMup() error {
	err := mh.repo.SaveMup(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// JWT middleware
func (mh *MupHandler) AuthorizeRoles(allowedRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
			tokenString := mh.extractTokenFromHeader(rr)
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

func (mh *MupHandler) extractTokenFromHeader(rr *http.Request) string {
	token := rr.Header.Get("Authorization")
	if token != "" {
		return token[len("Bearer "):]
	}
	return ""
}

func (mh *MupHandler) getJMBGFromToken(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	jmbg, ok1 := claims["sub"].(string)
	_, ok2 := claims["role"].(string)
	if !ok1 || !ok2 {
		return "", err
	}

	return jmbg, nil
}
