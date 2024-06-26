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
	sso   clients.SSOClient
}

// Constructor
func NewPoliceHandler(r *data.PoliceRepo, c clients.CourtClient, m clients.MupClient, s clients.SSOClient) *PoliceHandler {
	return &PoliceHandler{r, c, m, s}
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

	if driverCheck.AlcoholLevel < 0 {
		http.Error(w, "Alcohol level must be bigger than 0.", http.StatusBadRequest)
		return
	} else if driverCheck.AlcoholLevel > 0.2 {
		violation.Reason = fmt.Sprintf("drunk driving: %.2f \n", driverCheck.AlcoholLevel)
		violation.Description = "Driver was caught operating a vehicle with a blood alcohol level above the legal limit. \n"
	} else {
		log.Printf("Driver was caught operating a vehicle with a blood alcohol level within the legal limit.")
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
		log.Print("The driver is not under a driving ban.")
	}

	permit, err := ph.mup.GetDrivingPermitByJMBG(r.Context(), jmbgRequest, token)
	if err != nil {
		log.Printf("Failed to check driving permit: %v\n", err)
		http.Error(w, "Failed to check driving permit", http.StatusBadRequest)
		return
	}

	if permit.Number == "" {
		fmt.Printf("Not found driving permit.")
		http.Error(w, "Not found driving permit.", http.StatusBadRequest)
		return
	}

	if permit.ExpirationDate.Before(time.Now()) {
		violation.Reason += "Driving permit expired \n"
		violation.Description += "Driver was found to have an expired driving permit. \n"
		log.Print("Driving permit is expired")
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

	response := data.Response{}

	if alcoholLevel.AlcoholLevel > 0.2 {
		violation.Reason = fmt.Sprintf("drunk driving: %.2f \n", alcoholLevel.AlcoholLevel)
		violation.Description = "Driver was caught operating a vehicle with a blood alcohol level above the legal limit. \n"
	} else {
		response.Message = "Driver was caught operating a vehicle with a blood alcohol level within the legal limit."
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	token := ph.extractTokenFromHeader(r)
	_, err = ph.sso.GetPersonByJMBG(r.Context(), alcoholLevel.JMBG, token)
	if err != nil {
		http.Error(w, "Error with services communication", http.StatusBadRequest)
		log.Printf("Error while communicating with SSO service: %s", err.Error())
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

	response.Message = "Driver has more alcohol in his blood than is allowed."
	response.Data = violation

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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

	response := data.Response{}

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

	if permit.Number == "" {
		log.Printf("Not found driving permit.")
		http.Error(w, "Not found driving permit.", http.StatusBadRequest)
		return
	}

	if permit.ExpirationDate.Before(time.Now()) {
		violation.Reason += "Driving permit expired \n"
		violation.Description += "Driver was found to have an expired driving permit. \n"
		log.Print("Driving permit is expired")
	} else {
		response.Message = "The driver permit is valid."
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
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

	response.Message = "Driver has an expired driving permit."
	response.Data = violation

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (ph *PoliceHandler) CheckVehicleTire(w http.ResponseWriter, r *http.Request) {
	var tireType data.VehicleTireCheck
	err := json.NewDecoder(r.Body).Decode(&tireType)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v\n", err)
		return
	}

	response := data.Response{}
	token := ph.extractTokenFromHeader(r)

	now := time.Now()
	year := now.Year()

	startWinterPeriod := time.Date(year, time.November, 1, 0, 0, 0, 0, time.Local)
	endWinterPeriod := time.Date(year+1, time.April, 1, 0, 0, 0, 0, time.Local)
	startSummerPeriod := time.Date(year, time.April, 1, 0, 0, 0, 0, time.Local)
	endSummerPeriod := time.Date(year, time.November, 1, 0, 0, 0, 0, time.Local)

	if now.Month() >= time.November || now.Month() < time.April {
		endSummerPeriod = time.Date(year+1, time.November, 1, 0, 0, 0, 0, time.Local)
	}

	violation := data.TrafficViolation{
		ID:           primitive.NewObjectID(),
		Time:         now,
		ViolatorJMBG: tireType.JMBG,
		Location:     tireType.Location,
	}

	switch tireType.TireType {
	case "WINTER":
		if now.After(startWinterPeriod) && now.Before(endWinterPeriod) {
			response.Message = "No violation for winter tires"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
		violation.Reason = "Improper tire usage: WINTER tires outside summer period"
		violation.Description = "Driver was caught operating a vehicle with WINTER tires outside the summer period (April 1 to November 1), which is against regulations."

	case "SUMMER":
		if now.After(startSummerPeriod) && now.Before(endSummerPeriod) {
			response.Message = "No violation for summer tires in the summer period"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
		violation.Reason = "Improper tire usage: SUMMER tires during winter period"
		violation.Description = "Driver was caught operating a vehicle with SUMMER tires during the winter period (November 1 to April 1), which is against regulations."

	default:
		response.Message = "Invalid tire type specified"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = ph.sso.GetPersonByJMBG(r.Context(), tireType.JMBG, token)
	if err != nil {
		http.Error(w, "Error with services communication", http.StatusBadRequest)
		log.Printf("Error while communicating with SSO service: %s", err.Error())
		return
	}

	err = ph.repo.CreateTrafficViolation(r.Context(), &violation)
	if err != nil {
		response.Message = "Failed to create traffic violation"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		log.Printf("Failed to create traffic violation: %v\n", err)
		return
	}

	err = ph.court.CreateCrimeReport(r.Context(), violation, token)
	if err != nil {
		response.Message = "Failed to send crime report"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		log.Printf("Failed to send crime report: %v\n", err)
		return
	}

	response.Message = "Traffic violation created successfully."
	response.Data = violation
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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

	plates := data.PlateRequest{
		Plate: checkVehicleRegistration.PlatesNumber,
	}

	registration, err := ph.mup.GetRegistrationByPlate(r.Context(), plates, token)
	if err != nil {
		log.Printf("Failed to check registration by plate: %v\n", err)
		http.Error(w, "Failed to check registration by plate", http.StatusBadRequest)
		return
	}

	if registration.RegistrationNumber == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Wrong plates number in body request")
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
		log.Print("Vehicle registration has nor expired")
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
