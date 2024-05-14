package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"statistics/clients"
	"statistics/data"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ApplicationJson = "application/json"
const ContentType = "Content-Type"
const FailedToEncodeStatistics = "Failed to encode statistics"
const InvalidID = "Invalid ID"
const FailedToDecodeRequestBody = "Failed to decode request body"

var secretKey = []byte("eUpravaT2")

type StatisticsHandler struct {
	logger *log.Logger
	repo   *data.StatisticsRepo
	mup    clients.MupClient
}

func NewStatisticsHandler(l *log.Logger, r *data.StatisticsRepo, mc clients.MupClient) *StatisticsHandler {
	return &StatisticsHandler{l, r, mc}
}

// TODO Handler methods

func (sh *StatisticsHandler) GetAllTrafficStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accommodations, err := sh.repo.GetAllTrafficStatisticsData(ctx)
	if err != nil {
		sh.logger.Println("Failed to retrieve all traffic statistics")
		http.Error(rw, "Failed to retrieve traffic statistics", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(accommodations); err != nil {
		sh.logger.Println("Failed to encode all traffic statistics")
		http.Error(rw, FailedToEncodeStatistics, http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) GetAllCrimeStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accommodations, err := sh.repo.GetAllCrimeStatisticsData(ctx)
	if err != nil {
		sh.logger.Println("Failed to retrieve all crime statistics")
		http.Error(rw, "Failed to retrieve crime statistics", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(accommodations); err != nil {
		sh.logger.Println("Failed to encode all crime statistics")
		http.Error(rw, FailedToEncodeStatistics, http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) GetTrafficStatistic(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	sh.logger.Printf("Trying to retrieve traffic statistic with ID: %s\n", idStr)

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		http.Error(rw, InvalidID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	statistic, err := sh.repo.GetTrafficStatistic(ctx, id)
	if err != nil {
		sh.logger.Printf("Failed to retrieve traffic statistic with ID: %s, Error: %s\n", idStr, err.Error())
		http.Error(rw, "Failed to retrieve traffic statistic", http.StatusInternalServerError)
		return
	}

	if statistic == nil {
		sh.logger.Printf("Traffic statistic with ID: %s not found\n", idStr)
		http.NotFound(rw, r)
		return
	}

	sh.logger.Printf("Successfully retrieved traffic statistic with ID: %s\n", idStr)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(statistic); err != nil {
		sh.logger.Printf("Failed to encode traffic statistic with ID: %s, Error: %s\n", idStr, err.Error())
		http.Error(rw, "Failed to encode traffic statistic", http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) GetCrimeStatistic(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, InvalidID, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	statistic, err := sh.repo.GetCrimeStatistic(ctx, id)
	if err != nil {
		http.Error(rw, "Failed to retrieve crime statistic", http.StatusInternalServerError)
		return
	}

	if statistic == nil {
		http.NotFound(rw, r)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(statistic); err != nil {
		http.Error(rw, "Failed to encode crime statistic", http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) CreateTrafficStatistic(rw http.ResponseWriter, r *http.Request) {
	var statistic data.TrafficData
	if err := json.NewDecoder(r.Body).Decode(&statistic); err != nil {
		http.Error(rw, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	statistic.ID = primitive.NewObjectID()
	err := sh.repo.CreateTrafficStatisticData(r.Context(), &statistic)
	if err != nil {
		sh.logger.Println("Failed to create traffic statistic:", err)
		http.Error(rw, "Failed to create traffic statistic", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(statistic); err != nil {
		sh.logger.Println("Failed to encode traffic statistic:", err)
		http.Error(rw, FailedToEncodeStatistics, http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) CreateCrimeStatistic(rw http.ResponseWriter, r *http.Request) {
	var statistic data.CrimeData
	if err := json.NewDecoder(r.Body).Decode(&statistic); err != nil {
		http.Error(rw, "Failed to decode request body", http.StatusBadRequest)
		return
	}
	statistic.ID = primitive.NewObjectID()
	err := sh.repo.CreateCrimeStatisticData(r.Context(), &statistic)
	if err != nil {
		sh.logger.Println("Failed to create crime statistic:", err)
		http.Error(rw, "Failed to create crime statistic", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(statistic); err != nil {
		sh.logger.Println("Failed to encode crime statistic:", err)
		http.Error(rw, FailedToEncodeStatistics, http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) UpdateTrafficStatistic(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, InvalidID, http.StatusBadRequest)
		return
	}

	var updatedStatistic data.TrafficData
	if err := json.NewDecoder(r.Body).Decode(&updatedStatistic); err != nil {
		sh.logger.Println("Failed to decode request body:", err)
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		return
	}

	updatedStatistic.ID = id
	if err := sh.repo.UpdateTrafficStatistic(r.Context(), &updatedStatistic); err != nil {
		sh.logger.Println("Failed to update traffic statistic:", err)
		http.Error(rw, "Failed to update traffic statistic", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(updatedStatistic); err != nil {
		sh.logger.Println("Failed to encode updated traffic statistic:", err)
		http.Error(rw, "Failed to encode updated traffic statistic", http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) UpdateCrimeStatistic(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, InvalidID, http.StatusBadRequest)
		return
	}

	var updatedStatistic data.CrimeData
	if err := json.NewDecoder(r.Body).Decode(&updatedStatistic); err != nil {
		sh.logger.Println("Failed to decode request body:", err)
		http.Error(rw, FailedToDecodeRequestBody, http.StatusBadRequest)
		return
	}

	updatedStatistic.ID = id
	if err := sh.repo.UpdateCrimeStatistic(r.Context(), &updatedStatistic); err != nil {
		sh.logger.Println("Failed to update crime statistic:", err)
		http.Error(rw, "Failed to update crime statistic", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(updatedStatistic); err != nil {
		sh.logger.Println("Failed to encode updated crime statistic:", err)
		http.Error(rw, "Failed to encode updated crime statistic", http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) DeleteTrafficStatistic(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, InvalidID, http.StatusBadRequest)
		return
	}

	if err := sh.repo.DeleteTrafficStatistic(r.Context(), id); err != nil {
		sh.logger.Println("Failed to delete traffic statistic:", err)
		http.Error(rw, "Failed to delete traffic statistic", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (sh *StatisticsHandler) DeleteCrimeStatistic(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(vars["id"])
	if err != nil {
		http.Error(rw, InvalidID, http.StatusBadRequest)
		return
	}

	if err := sh.repo.DeleteCrimeStatistic(r.Context(), id); err != nil {
		sh.logger.Println("Failed to delete crime statistic:", err)
		http.Error(rw, "Failed to delete crime statistic", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (sh *StatisticsHandler) GetVehicleStatisticsByYear(rw http.ResponseWriter, r *http.Request) {
	vehicles, err := sh.mup.GetAllRegisteredVehicles(r.Context())
	if err != nil {
		sh.logger.Println("Failed to retrieve vehicles:", err)
		http.Error(rw, "Failed to retrieve vehicles", http.StatusInternalServerError)
		return
	}

	vehicleStatistics := make(map[int]int)

	for _, vehicle := range vehicles {
		year := vehicle.Year
		if _, ok := vehicleStatistics[year]; ok {
			vehicleStatistics[year]++
		} else {
			vehicleStatistics[year] = 1
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(vehicleStatistics); err != nil {
		sh.logger.Println("Failed to encode vehicle statistics:", err)
		http.Error(rw, "Failed to encode vehicle statistics", http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) GetRegisteredVehicles(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicles, err := sh.mup.GetAllRegisteredVehicles(ctx)
	if err != nil {
		sh.logger.Println("Failed to retrieve registered vehicles:", err)
		http.Error(rw, "Failed to retrieve registered vehicles", http.StatusInternalServerError)
		return
	}

	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(vehicles); err != nil {
		sh.logger.Println("Failed to encode registered vehicles:", err)
		http.Error(rw, "Failed to encode registered vehicles", http.StatusInternalServerError)
	}
}

func (sh *StatisticsHandler) GetMostPopularBrands(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicles, err := sh.mup.GetAllRegisteredVehicles(ctx)
	if err != nil {
		sh.logger.Println("Failed to retrieve registered vehicles:", err)
		http.Error(rw, "Failed to retrieve registered vehicles", http.StatusInternalServerError)
		return
	}

	// Izračunavanje broja vozila po brendovima
	brandCount := make(map[string]int)

	for _, vehicle := range vehicles {
		brandCount[vehicle.Brand]++
	}

	// Prikaz rezultata
	rw.Header().Set(ContentType, ApplicationJson)
	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(brandCount); err != nil {
		sh.logger.Println("Failed to encode most popular brands:", err)
		http.Error(rw, "Failed to encode most popular brands", http.StatusInternalServerError)
	}
}

// JWT middleware
func (sh *StatisticsHandler) AuthorizeRoles(allowedRoles ...string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, rr *http.Request) {
			tokenString := sh.extractTokenFromHeader(rr)
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

func (sh *StatisticsHandler) extractTokenFromHeader(rr *http.Request) string {
	token := rr.Header.Get("Authorization")
	if token != "" {
		return token[len("Bearer "):]
	}
	return ""
}
