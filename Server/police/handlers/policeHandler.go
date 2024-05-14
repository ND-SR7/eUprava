package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"police/clients"
	"police/data"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeyProduct struct{}

var secretKey = []byte("eUpravaT2")

type PoliceHandler struct {
	repo  *data.PoliceRepo
	court clients.CourtClient
}

// Constructor
func NewPoliceHandler(r *data.PoliceRepo, c clients.CourtClient) *PoliceHandler {
	return &PoliceHandler{r, c}
}

// Ping
func (ph *PoliceHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Pong"))

	w.WriteHeader(http.StatusOK)
}

func (ph *PoliceHandler) CreateTrafficViolation(w http.ResponseWriter, r *http.Request) {
	var violation data.TrafficViolation
	err := json.NewDecoder(r.Body).Decode(&violation)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
	if err != nil {
		http.Error(w, "Failed to create traffic violation", http.StatusInternalServerError)
		log.Printf("Failed to create traffic violation: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(violation)
}

func (ph *PoliceHandler) CheckAlcoholLevel(w http.ResponseWriter, r *http.Request) {
	var alcoholTest data.AlcoholTest
	err := json.NewDecoder(r.Body).Decode(&alcoholTest)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	if alcoholTest.AlcoholLevel > 0.2 {
		violation := data.TrafficViolation{
			Reason:        fmt.Sprintf("drunk driving: %.2f", alcoholTest.AlcoholLevel),
			Description:   "Driver was caught operating a vehicle with a blood alcohol level above the legal limit.",
			Time:          time.Now(),
			ViolatorEmail: alcoholTest.UserEmail,
			Location:      alcoholTest.Location,
		}

		err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
		if err != nil {
			http.Error(w, "Failed to create traffic violation", http.StatusInternalServerError)
			log.Printf("Failed to create traffic violation: %v\n", err)
			return
		}

		token := ph.extractTokenFromHeader(r)

		err := ph.court.CreateCrimeReport(r.Context(), violation, token)
		if err != nil {
			http.Error(w, "Failed to send crime report", http.StatusInternalServerError)
			log.Printf("Failed to send crime report: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(violation)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Driver is not under the influence of alcohol"))
}

func (ph *PoliceHandler) GetTrafficViolationByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	violationID, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing violation ID", http.StatusBadRequest)
		log.Println("Missing violation ID in request")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(violationID)
	if err != nil {
		http.Error(w, "Invalid violation ID", http.StatusBadRequest)
		log.Printf("Invalid violation ID: %v\n", err)
		return
	}

	violation, err := ph.repo.GetTrafficViolationByID(r.Context(), objectID)
	if err != nil {
		http.Error(w, "Failed to retrieve traffic violation", http.StatusInternalServerError)
		log.Printf("Failed to retrieve traffic violation: %v\n", err)
		return
	}

	json.NewEncoder(w).Encode(violation)
}

func (ph *PoliceHandler) GetAllTrafficViolations(w http.ResponseWriter, r *http.Request) {
	violations, err := ph.repo.GetAllTrafficViolations(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve traffic violations", http.StatusInternalServerError)
		log.Printf("Failed to retrieve traffic violations: %v\n", err)
		return
	}

	json.NewEncoder(w).Encode(violations)
}

func (ph *PoliceHandler) UpdateTrafficViolation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	violationID, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing violation ID", http.StatusBadRequest)
		log.Println("Missing violation ID in request")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(violationID)
	if err != nil {
		http.Error(w, "Invalid violation ID", http.StatusBadRequest)
		log.Printf("Invalid violation ID: %v\n", err)
		return
	}

	existingViolation, err := ph.repo.GetTrafficViolationByID(r.Context(), objectID)
	if err != nil {
		http.Error(w, "Failed to retrieve traffic violation", http.StatusInternalServerError)
		log.Printf("Failed to retrieve traffic violation: %v\n", err)
		return
	}

	var update data.TrafficViolation
	err = json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	if update.ViolatorEmail != "" {
		existingViolation.ViolatorEmail = update.ViolatorEmail
	}
	if update.Reason != "" {
		existingViolation.Reason = update.Reason
	}
	if update.Description != "" {
		existingViolation.Description = update.Description
	}
	if !update.Time.IsZero() {
		existingViolation.Time = update.Time
	}
	if update.Location != "" {
		existingViolation.Location = update.Location
	}

	err = ph.repo.UpdateTrafficViolation(r.Context(), objectID, existingViolation)
	if err != nil {
		http.Error(w, "Failed to update traffic violation", http.StatusInternalServerError)
		log.Printf("Failed to update traffic violation: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingViolation)
}

func (ph *PoliceHandler) DeleteTrafficViolation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	violationID, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing violation ID", http.StatusBadRequest)
		log.Println("Missing violation ID in request")
		return
	}

	objectID, err := primitive.ObjectIDFromHex(violationID)
	if err != nil {
		http.Error(w, "Invalid violation ID", http.StatusBadRequest)
		log.Printf("Invalid violation ID: %v\n", err)
		return
	}

	err = ph.repo.DeleteTrafficViolation(r.Context(), objectID)
	if err != nil {
		http.Error(w, "Failed to delete traffic violation", http.StatusInternalServerError)
		log.Printf("Failed to delete traffic violation: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Traffic violation deleted successfully"})
}

// JWT middleware
func (ph *PoliceHandler) AuthorizeRoles(allowedRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
			tokenString := ph.extractTokenFromHeader(rr)
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

func (ph *PoliceHandler) extractTokenFromHeader(rr *http.Request) string {
	token := rr.Header.Get("Authorization")
	if token != "" {
		return token[len("Bearer "):]
	}
	return ""
}
