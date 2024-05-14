package main

import (
	"context"
	"log"
	"net/http"
	"statistics/clients"
	"statistics/data"

	"os"
	"os/signal"
	"statistics/handlers"
	"syscall"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const TrafficStatisticPath = "/api/v1/traffic-statistic/{id}"
const CrimeStatisticPath = "/api/v1/crime-statistic/{id}"

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8084"
	}

	// Context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Logger init

	logger := log.New(os.Stdout, "[statistics-service] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[statistics-store] ", log.LstdFlags)

	// DB init & ping

	store, err := data.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Disconnect(timeoutContext)
	store.Ping()

	mupClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	mup := clients.NewMupClient(mupClient, os.Getenv("MUP_SERVICE_URI"))

	// TODO: Handler init

	statisticsHandler := handlers.NewStatisticsHandler(logger, store, mup)

	router := mux.NewRouter()

	// TODO: Router methods

	createTrafficStatisticRouter := router.Methods(http.MethodPost).Path("/api/v1/traffic-statistic").Subrouter()
	createTrafficStatisticRouter.HandleFunc("", statisticsHandler.CreateTrafficStatistic)

	createCrimeStatisticRouter := router.Methods(http.MethodPost).Path("/api/v1/crime-statistic").Subrouter()
	createCrimeStatisticRouter.HandleFunc("", statisticsHandler.CreateTrafficStatistic)

	getTrafficStatisticRouter := router.Methods(http.MethodGet).Path(TrafficStatisticPath).Subrouter()
	getTrafficStatisticRouter.HandleFunc("", statisticsHandler.GetTrafficStatistic)

	getCrimeStatisticRouter := router.Methods(http.MethodGet).Path(CrimeStatisticPath).Subrouter()
	getCrimeStatisticRouter.HandleFunc("", statisticsHandler.GetCrimeStatistic)

	getAllTrafficStatisticRouter := router.Methods(http.MethodGet).Path("/api/v1/traffic-statistic").Subrouter()
	getAllTrafficStatisticRouter.HandleFunc("", statisticsHandler.GetAllTrafficStatistics)

	getAllCrimeStatisticRouter := router.Methods(http.MethodGet).Path("/api/v1/crime-statistic").Subrouter()
	getAllCrimeStatisticRouter.HandleFunc("", statisticsHandler.GetAllCrimeStatistics)

	updateTrafficStatisticRouter := router.Methods(http.MethodPut).Path(TrafficStatisticPath).Subrouter()
	updateTrafficStatisticRouter.HandleFunc("", statisticsHandler.UpdateTrafficStatistic)

	updateCrimeStatisticRouter := router.Methods(http.MethodPut).Path(CrimeStatisticPath).Subrouter()
	updateCrimeStatisticRouter.HandleFunc("", statisticsHandler.UpdateCrimeStatistic)

	deleteTrafficStatisticRouter := router.Methods(http.MethodDelete).Path(TrafficStatisticPath).Subrouter()
	deleteTrafficStatisticRouter.HandleFunc("", statisticsHandler.DeleteTrafficStatistic)

	deleteCrimeStatisticRouter := router.Methods(http.MethodDelete).Path(CrimeStatisticPath).Subrouter()
	deleteCrimeStatisticRouter.HandleFunc("", statisticsHandler.DeleteCrimeStatistic)

	getVehicleStatisticsByYearRouter := router.Methods(http.MethodGet).Path("/api/v1/vehicle-statistics-by-year").Subrouter()
	getVehicleStatisticsByYearRouter.HandleFunc("", statisticsHandler.GetVehicleStatisticsByYear)

	getRegisteredVehiclesRouter := router.Methods(http.MethodGet).Path("/api/v1/registered-vehicles").Subrouter()
	getRegisteredVehiclesRouter.HandleFunc("", statisticsHandler.GetRegisteredVehicles)

	getMostPopularVehicleBrendsRouter := router.Methods(http.MethodGet).Path("/api/v1/most-popular-vehicles").Subrouter()
	getMostPopularVehicleBrendsRouter.HandleFunc("", statisticsHandler.GetMostPopularBrands)

	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type"}),
	)

	pingRouter := router.Methods("GET").Subrouter()
	pingRouter.HandleFunc("/api/v1", statisticsHandler.Ping).Methods("GET")
	pingRouter.Use(cors)
	pingRouter.Use(statisticsHandler.AuthorizeRoles("USER", "ADMIN"))

	// Initialize the server
	server := http.Server{
		Addr:         ":" + port,
		Handler:      cors(router),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	logger.Printf("Server listening on port: %s\n", port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatalf("Error while serving request: %v\n", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGTERM)

	sig := <-sigCh
	logger.Printf("Recieved terminate, starting gracefull shutdown: %v\n", sig)

	// Gracefull shutdown
	if server.Shutdown(timeoutContext) != nil {
		logger.Fatalln("Cannot gracefully shutdown")
	}
	logger.Println("Server gracefully stopped")
}
