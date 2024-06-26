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
	mup   clients.MupClient
}

// Constructor
func NewPoliceHandler(r *data.PoliceRepo, c clients.CourtClient, m clients.MupClient) *PoliceHandler {
	return &PoliceHandler{r, c, m}
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

	violation.ID = primitive.NewObjectID()

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

func (ph *PoliceHandler) CheckAll(w http.ResponseWriter, r *http.Request) {
	var driverCheck data.DriverCheck
	err := json.NewDecoder(r.Body).Decode(&driverCheck)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	violation := data.TrafficViolation{
		ID:           primitive.NewObjectID(),
		Time:         time.Now(),
		ViolatorJMBG: driverCheck.JMBG,
		Location:     driverCheck.Location,
	}

	if driverCheck.AlcoholLevel <= 0 {
		violation.Reason = "Driver is not under the influence of alcohol \n"
		violation.Description = "Driver was caught operating a vehicle with a blood alcohol level within the legal limit. \n"
	} else {

		if driverCheck.AlcoholLevel > 0.2 {
			violation.Reason = fmt.Sprintf("drunk driving: %.2f \n", driverCheck.AlcoholLevel)
			violation.Description = "Driver was caught operating a vehicle with a blood alcohol level above the legal limit. \n"
		} else {
			violation.Reason = fmt.Sprintf("drunk driving: %.2f", driverCheck.AlcoholLevel)
			violation.Description = "Driver was caught operating a vehicle with a blood alcohol level within the legal limit."
		}
	}

	if driverCheck.Tire != "" {
		switch driverCheck.Tire {
		case "SUMMER":
			now := time.Now()
			month := now.Month()
			day := now.Day()

			winterTiresRequired := (month >= time.November && month <= time.December) || (month == time.January && day <= 1)

			if winterTiresRequired {
				violation.Reason += fmt.Sprintf("Winter tires required, provided tire type: %s \n", driverCheck.Tire)
				violation.Description += "Additionally, the vehicle was equipped with incorrect tires for the current date."
			}

		case "WINTER":
			now := time.Now()
			month := now.Month()
			day := now.Day()

			winterTiresRequired := (month >= time.November && month <= time.December) || (month == time.January && day <= 1)
			if !winterTiresRequired {
				violation.Reason += fmt.Sprintf("Winter tires not required, provided tire type: %s \n", driverCheck.Tire)
				violation.Description += "The vehicle was equipped with winter tires outside of the required period. \n"
			}

		default:
			http.Error(w, "Invalid tire type provided. Must be either SUMMER or WINTER", http.StatusBadRequest)
			return
		}

	}

	token := ph.extractTokenFromHeader(r)

	if driverCheck.JMBG != "" {
		jmbgRequest := data.JMBGRequest{JMBG: driverCheck.JMBG}
		drivingBan, err := ph.mup.CheckDrivigBan(r.Context(), jmbgRequest, token)
		if err != nil {
			http.Error(w, "Failed to check driving ban: "+err.Error(), http.StatusBadRequest)
			log.Printf("Failed to check driving ban: %v\n", err)
			return
		}

		if drivingBan {
			violation.Reason += "Driving ban is effect \n"
			violation.Description += "Driver was found to be operating a vehicle under active driving ban. \n"
			log.Print("drivingBan is true")
		} else {
			violation.Reason += "No driving ban \n"
			violation.Description += "Driver was not found to be operating a vehicle under any driving ban. \n"
			log.Print("drivingBan is false")
		}

		permit, err := ph.mup.GetDrivingPermitByJMBG(r.Context(), jmbgRequest, token)
		if err != nil {
			log.Printf("Failed to check driving permit: %v\n", err)
			http.Error(w, "Failed to check driving permit", http.StatusBadRequest)
			return
		}

		if permit.ExpirationDate.Before(time.Now()) {
			violation.Reason += "Driving permit expired \n"
			violation.Description += "Driver was found to have an expired driving permit. \n"
			log.Print("Driving permit is expired")
		}

	}

	err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
	if err != nil {
		http.Error(w, "Failed to create traffic violation", http.StatusInternalServerError)
		log.Printf("Failed to create traffic violation: %v\n", err)
		return
	}

	err = ph.court.CreateCrimeReport(r.Context(), violation, token)
	if err != nil {
		http.Error(w, "Failed to send crime report", http.StatusInternalServerError)
		log.Printf("Failed to send crime report: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(violation)
}

func (ph *PoliceHandler) CheckAlcoholLevel(w http.ResponseWriter, r *http.Request) {
	var alcoholLevel data.AlcoholRequest

	err := json.NewDecoder(r.Body).Decode(&alcoholLevel)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	violation := data.TrafficViolation{
		ID:           primitive.NewObjectID(),
		Time:         time.Now(),
		ViolatorJMBG: alcoholLevel.JMBG,
		Location:     alcoholLevel.Location,
	}

	if alcoholLevel.AlcoholLevel > 0.2 {
		violation.Reason = fmt.Sprintf("drunk driving: %.2f \n", alcoholLevel.AlcoholLevel)
		violation.Description = "Driver was caught operating a vehicle with a blood alcohol level above the legal limit. \n"
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Driver was caught operating a vehicle with a blood alcohol level within the legal limit.")
		return
	}

	token := ph.extractTokenFromHeader(r)

	err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
	if err != nil {
		http.Error(w, "Failed to create traffic violation", http.StatusBadRequest)
		log.Printf("Failed to create traffic violation: %v\n", err)
		return
	}

	err = ph.court.CreateCrimeReport(r.Context(), violation, token)
	if err != nil {
		http.Error(w, "Failed to send crime report", http.StatusBadRequest)
		log.Printf("Failed to send crime report: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(violation)
}

func (ph *PoliceHandler) CheckDriverBan(w http.ResponseWriter, r *http.Request) {
	var driverBan data.DriverBanAndPermitRequest
	err := json.NewDecoder(r.Body).Decode(&driverBan)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	violation := data.TrafficViolation{
		ID:           primitive.NewObjectID(),
		Time:         time.Now(),
		ViolatorJMBG: driverBan.JMBG,
		Location:     driverBan.Location,
	}

	token := ph.extractTokenFromHeader(r)

	jmbgRequest := data.JMBGRequest{
		JMBG: driverBan.JMBG,
	}

	drivingBan, err := ph.mup.CheckDrivigBan(r.Context(), jmbgRequest, token)
	if err != nil {
		http.Error(w, "Failed to check driving ban: "+err.Error(), http.StatusBadRequest)
		log.Printf("Failed to check driving ban: %v\n", err)
		return
	}

	if drivingBan {
		violation.Reason += "Driving ban is effect \n"
		violation.Description += "Driver was found to be operating a vehicle under active driving ban. \n"
		log.Print("drivingBan is true")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "The driver is not under a driving ban.")
		return
	}

	err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
	if err != nil {
		http.Error(w, "Failed to create traffic violation", http.StatusInternalServerError)
		log.Printf("Failed to create traffic violation: %v\n", err)
		return
	}

	err = ph.court.CreateCrimeReport(r.Context(), violation, token)
	if err != nil {
		http.Error(w, "Failed to send crime report", http.StatusInternalServerError)
		log.Printf("Failed to send crime report: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(violation)
}

func (ph *PoliceHandler) CheckDriverPermitValidity(w http.ResponseWriter, r *http.Request) {
	var driverBan data.DriverBanAndPermitRequest
	err := json.NewDecoder(r.Body).Decode(&driverBan)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	violation := data.TrafficViolation{
		ID:           primitive.NewObjectID(),
		Time:         time.Now(),
		ViolatorJMBG: driverBan.JMBG,
		Location:     driverBan.Location,
	}

	token := ph.extractTokenFromHeader(r)

	jmbgRequest := data.JMBGRequest{
		JMBG: driverBan.JMBG,
	}

	permit, err := ph.mup.GetDrivingPermitByJMBG(r.Context(), jmbgRequest, token)
	if err != nil {
		log.Printf("Failed to check driving permit: %v\n", err)
		http.Error(w, "Failed to check driving permit", http.StatusBadRequest)
		return
	}

	if permit.ExpirationDate.Before(time.Now()) {
		violation.Reason += "Driving permit expired \n"
		violation.Description += "Driver was found to have an expired driving permit. \n"
		log.Print("Driving permit is expired")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "The driver permit is valid.")
		return
	}

	err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
	if err != nil {
		http.Error(w, "Failed to create traffic violation", http.StatusInternalServerError)
		log.Printf("Failed to create traffic violation: %v\n", err)
		return
	}

	err = ph.court.CreateCrimeReport(r.Context(), violation, token)
	if err != nil {
		http.Error(w, "Failed to send crime report", http.StatusInternalServerError)
		log.Printf("Failed to send crime report: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(violation)
}

func (ph *PoliceHandler) CheckVehicleTire(w http.ResponseWriter, r *http.Request) {
	var tireType data.VehicleTireCheck
	err := json.NewDecoder(r.Body).Decode(&tireType)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	violation := data.TrafficViolation{
		ID:           primitive.NewObjectID(),
		Time:         time.Now(),
		ViolatorJMBG: tireType.JMBG,
		Location:     tireType.Location,
	}

	token := ph.extractTokenFromHeader(r)

	// currentDate := time.Now()
	startWinterPeriod := time.Date((time.Now()).Year(), time.November, 1, 0, 0, 0, 0, time.Local)
	endWinterPeriod := time.Date((time.Now()).Year(), time.April, 1, 0, 0, 0, 0, time.Local)

	if tireType.TireType == "WINTER" {
		// No violation for winter tires
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "No violation for winter tires")
		return
	} else if tireType.TireType == "SUMMER" && (time.Now().After(startWinterPeriod) || time.Now().Before(endWinterPeriod)) {
		violation.Reason = "Improper tire usage: SUMMER tires during winter period"
		violation.Description = "Driver was caught operating a vehicle with SUMMER tires during the winter period (November 1 to April 1), which is against regulations."

		err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
		if err != nil {
			http.Error(w, "Failed to create traffic violation", http.StatusBadRequest)
			log.Printf("Failed to create traffic violation: %v\n", err)
			return
		}

		err = ph.court.CreateCrimeReport(r.Context(), violation, token)
		if err != nil {
			http.Error(w, "Failed to send crime report", http.StatusBadRequest)
			log.Printf("Failed to send crime report: %v\n", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(violation)
	} else {
		// No violation for summer tires outside the winter period
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "No violation for summer tires outside the winter period")
		return
	}
}

func (ph *PoliceHandler) CheckVehicleRegistration(w http.ResponseWriter, r *http.Request) {
	var checkVehicleRegistration data.CheckVehicleRegistration
	err := json.NewDecoder(r.Body).Decode(&checkVehicleRegistration)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	violation := data.TrafficViolation{
		ID:           primitive.NewObjectID(),
		Time:         time.Now(),
		ViolatorJMBG: checkVehicleRegistration.JMBG,
		Location:     checkVehicleRegistration.Location,
	}

	token := ph.extractTokenFromHeader(r)

	registration, err := ph.mup.GetVehicleRegistration(r.Context(), checkVehicleRegistration, token)
	if err != nil {
		log.Printf("Failed to check driving permit: %v\n", err)
		http.Error(w, "Failed to check driving permit", http.StatusBadRequest)
		return
	}

	if registration.ExpirationDate.Before(time.Now()) {
		violation.Reason = "Vehicle registration expired"
		violation.Description = "Driver was found to be operating a vehicle with an expired registration."
		log.Print("Vehicle registration is expired")
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "The vehicle registration has not expired.")
		return
	}

	err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
	if err != nil {
		http.Error(w, "Failed to create traffic violation", http.StatusBadRequest)
		log.Printf("Failed to create traffic violation: %v\n", err)
		return
	}

	err = ph.court.CreateCrimeReport(r.Context(), violation, token)
	if err != nil {
		http.Error(w, "Failed to send crime report", http.StatusBadRequest)
		log.Printf("Failed to send crime report: %v\n", err)
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

	if update.ViolatorJMBG != "" {
		existingViolation.ViolatorJMBG = update.ViolatorJMBG
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
