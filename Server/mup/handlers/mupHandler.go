package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"mup/data"
	"net/http"
)

const ApplicationJson = "application/json"
const ContentType = "Content-Type"
const FailedToEncodePerson = "Failed to encode person"
const FailedToEcnodeDrivingBans = "Failed to encode driving bans"
const InvalidID = "Invalid ID"
const FailedToDecodeRequestBody = "Failed to decode request body"

type KeyProduct struct{}

var secretKey = []byte("UpravaT2")

type MupHandler struct {
	repo   *data.MUPRepo
	logger *log.Logger
}

func NewMupHandler(r *data.MUPRepo, log *log.Logger) *MupHandler {
	return &MupHandler{r, log}
}

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

	if err := json.NewDecoder(r.Body).Decode(&trafficPermit); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	if err := mh.repo.SubmitTrafficPermitRequest(r.Context(), &trafficPermit); err != nil {
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
func (mh *MupHandler) SaveMup(mup data.Mup) error {
	err := mh.repo.SaveMup(context.Background(), mup)
	if err != nil {
		return err
	}
	return nil
}

// Middlerware
func (mh *MupHandler) MiddlewarePersonDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		patient := &data.Person{}
		err := patient.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			mh.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, patient)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
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
