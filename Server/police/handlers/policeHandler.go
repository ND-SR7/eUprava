package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"police/data"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PoliceHandler struct {
	repo *data.PoliceRepo
}

// Constructor
func NewPoliceHandler(r *data.PoliceRepo) *PoliceHandler {
	return &PoliceHandler{r}
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
