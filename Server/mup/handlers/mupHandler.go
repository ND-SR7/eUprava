package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"mup/data"
	"mup/services"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

const ApplicationJson = "application/json"
const ContentType = "Content-Type"
const FailedToReadUsernameFromToken = "Failed to read username from token"
const FailedToEncodeDrivingBans = "Failed to encode driving bans"
const FailedToDecodeRequestBody = "Failed to decode request body"

type KeyProduct struct{}

var secretKey = []byte("eUpravaT2")

type MupHandler struct {
	service *services.MupService
	logger  *log.Logger
}

func NewMupHandler(service *services.MupService, logger *log.Logger) *MupHandler {
	return &MupHandler{service: service, logger: logger}
}

// Ping
func (mh *MupHandler) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Pong"))

	w.WriteHeader(http.StatusOK)
}

// GET METHODS
func (mh *MupHandler) CheckForPersonsDrivingBans(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenStr := mh.extractTokenFromHeader(r)

	jmbg, err := mh.getJMBGFromToken(tokenStr)
	if err != nil {
		fmt.Printf("Error while reading JMBG from token: %v", err)
		http.Error(rw, FailedToReadUsernameFromToken, http.StatusBadRequest)
		return
	}

	drivingBans, err := mh.service.CheckForPersonsDrivingBans(ctx, jmbg)
	if err != nil {
		http.Error(rw, "Failed to retrieve persons driving bans", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(drivingBans); err != nil {
		http.Error(rw, FailedToEncodeDrivingBans, http.StatusInternalServerError)
	}
	fmt.Println("Successfully fetched driving bans")
}

func (mh *MupHandler) CheckForRegisteredVehicles(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicles, err := mh.service.RetrieveRegisteredVehicles(ctx)
	if err != nil {
		http.Error(rw, "Failed to retrieve registered vehicles", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(vehicles); err != nil {
		http.Error(rw, FailedToEncodeDrivingBans, http.StatusInternalServerError)
	}
	fmt.Println("Successfully fetched registered vehicles")
}

func (mh *MupHandler) GetPersonsRegistrations(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenStr := mh.extractTokenFromHeader(r)
	jmbg, err := mh.getJMBGFromToken(tokenStr)
	if err != nil {
		http.Error(rw, "Failed to read JMBG from token", http.StatusBadRequest)
		return
	}

	registrations, err := mh.service.GetPersonsRegistrations(ctx, jmbg)
	if err != nil {
		http.Error(rw, "Failed to retrieve user registrations", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(registrations); err != nil {
		http.Error(rw, "Failed to encode registrations", http.StatusInternalServerError)
	}
}

func (mh *MupHandler) GetUserDrivingPermit(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tokenStr := mh.extractTokenFromHeader(r)
	jmbg, err := mh.getJMBGFromToken(tokenStr)
	if err != nil {
		http.Error(rw, "Failed to read JMBG from token", http.StatusBadRequest)
		return
	}

	drivingPermits, err := mh.service.GetUserDrivingPermit(ctx, jmbg)
	if err != nil {
		http.Error(rw, "Failed to retrieve user driving permits", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(drivingPermits); err != nil {
		http.Error(rw, "Failed to encode driving permits", http.StatusInternalServerError)
	}
}

func (mh *MupHandler) GetPendingRegistrationRequests(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pendingRequests, err := mh.service.GetPendingRegistrationRequests(ctx)
	if err != nil {
		http.Error(rw, "Failed to retrieve pending registration requests", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(pendingRequests); err != nil {
		http.Error(rw, "Failed to encode pending requests", http.StatusInternalServerError)
	}
}

func (mh *MupHandler) GetDrivingBan(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request data.JMBGRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	drivingBan, err := mh.service.GetDrivingBan(ctx, request.JMBG)
	if err != nil {
		http.Error(rw, "Failed to retrieve driving ban", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(drivingBan); err != nil {
		http.Error(rw, "Failed to encode driving ban", http.StatusInternalServerError)
	}
}

func (mh *MupHandler) GetPendingTrafficPermitRequests(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pendingRequests, err := mh.service.GetPendingTrafficPermitRequests(ctx)
	if err != nil {
		http.Error(rw, "Failed to retrieve pending traffic permit requests", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(pendingRequests); err != nil {
		http.Error(rw, "Failed to encode pending requests", http.StatusInternalServerError)
	}
}

func (mh *MupHandler) GetRegistrationByPlate(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request data.PlateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	registration, err := mh.service.GetRegistrationByPlate(ctx, request.Plate)
	if err != nil {
		http.Error(rw, "Failed to retrieve registration", http.StatusInternalServerError)
		return
	}

	if registration.RegistrationNumber == "" {
		registration = data.Registration{}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(registration); err != nil {
		http.Error(rw, "Failed to encode registration", http.StatusInternalServerError)
	}
}

func (mh *MupHandler) GetDrivingPermitByJMBG(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request data.JMBGRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	drivingPermit, err := mh.service.GetDrivingPermitByJMBG(ctx, request.JMBG)
	if err != nil {
		http.Error(rw, "Failed to retrieve driving permit", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(drivingPermit); err != nil {
		http.Error(rw, "Failed to encode driving permit", http.StatusInternalServerError)
	}
}

func (mh *MupHandler) GetPersonsVehicles(rw http.ResponseWriter, r *http.Request) {
	tokenStr := mh.extractTokenFromHeader(r)
	jmbg, err := mh.getJMBGFromToken(tokenStr)
	if err != nil {
		fmt.Printf("Error while reading JMBG from token: %v", err)
		http.Error(rw, FailedToReadUsernameFromToken, http.StatusBadRequest)
		return
	}

	vehicles, err := mh.service.GetPersonsVehicles(r.Context(), jmbg)
	if err != nil {
		log.Printf("Failed to retrieve vehicles for person with JMBG %s: %v", jmbg, err)
		http.Error(rw, "Failed to retrieve vehicles", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(vehicles); err != nil {
		log.Printf("Failed to encode vehicles: %v", err)
		http.Error(rw, "Failed to encode vehicles", http.StatusInternalServerError)
	}
}

func (mh *MupHandler) GetVehiclesDTOByJMBG(rw http.ResponseWriter, r *http.Request) {
	tokenStr := mh.extractTokenFromHeader(r)
	jmbg, err := mh.getJMBGFromToken(tokenStr)
	if err != nil {
		fmt.Printf("Error while reading JMBG from token: %v", err)
		http.Error(rw, FailedToReadUsernameFromToken, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	vehicleDTOs, err := mh.service.GetVehiclesDTOByJMBG(ctx, jmbg)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(vehicleDTOs); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

// POST METHODS
func (mh *MupHandler) SubmitRegistrationRequest(rw http.ResponseWriter, r *http.Request) {
	var registration data.Registration

	tokenStr := mh.extractTokenFromHeader(r)
	jmbg, err := mh.getJMBGFromToken(tokenStr)
	if err != nil {
		http.Error(rw, "Failed to read JMBG from token", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&registration); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	registration.Owner = jmbg

	if err := mh.service.SubmitRegistrationRequest(r.Context(), &registration); err != nil {
		log.Printf("Failed to submit registration request: %v", err)
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
		fmt.Printf("Error while reading JMBG from token: %v", err)
		http.Error(rw, FailedToReadUsernameFromToken, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&trafficPermit); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	trafficPermit.Person = jmbg

	if err := mh.service.SubmitTrafficPermitRequest(ctx, &trafficPermit, jmbg, tokenStr); err != nil {
		log.Printf("Failed to submit traffic permit request: %v", err)
		http.Error(rw, "Failed to submit traffic permit request", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(trafficPermit); err != nil {
		log.Printf("Failed to encode traffic permit request: %v", err)
		http.Error(rw, "Failed to encode traffic permit request", http.StatusInternalServerError)
	}

	log.Printf("Successfully created traffic permit request with id '%s'", trafficPermit.ID.Hex())
}

func (mh *MupHandler) SaveVehicle(rw http.ResponseWriter, r *http.Request) {
	var vehicle data.Vehicle
	tokenStr := mh.extractTokenFromHeader(r)

	jmbg, err := mh.getJMBGFromToken(tokenStr)
	if err != nil {
		fmt.Printf("Error while reading JMBG from token: %v", err)
		http.Error(rw, FailedToReadUsernameFromToken, http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	vehicle.Owner = jmbg

	if err := mh.service.SaveVehicle(r.Context(), &vehicle); err != nil {
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

	if err := mh.service.IssueDrivingBan(r.Context(), &drivingBan); err != nil {
		log.Printf("Failed to issue driving ban: %v", err)
		http.Error(rw, "Failed to issue driving ban", http.StatusInternalServerError)
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

	if err := mh.service.ApproveRegistration(r.Context(), registration); err != nil {
		log.Printf("Failed to approve registration: %v", err)
		http.Error(rw, "Failed to approve registration", http.StatusInternalServerError)
		return
	}

	registration.Approved = true

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

	trafficPermit.Approved = true

	if err := mh.service.ApproveTrafficPermitRequest(r.Context(), trafficPermit.ID); err != nil {
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

// DELETE METHODS

func (mh *MupHandler) DeletePendingRegistration(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	request := vars["request"]

	if request == "" {
		http.Error(rw, "registration number is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err := mh.service.DeletePendingRegistration(ctx, request)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (mh *MupHandler) DeletePendingTrafficPermit(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	request := vars["request"]

	if request == "" {
		http.Error(rw, "permit ID is required", http.StatusBadRequest)
		return
	}

	permitID, err := primitive.ObjectIDFromHex(request)
	if err != nil {
		http.Error(rw, "Invalid permit ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = mh.service.DeletePendingTrafficPermit(ctx, permitID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
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
