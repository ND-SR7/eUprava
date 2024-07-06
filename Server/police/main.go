package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"police/clients"
	"police/data"
	"police/handlers"
	"syscall"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8082"
	}

	// Context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Logger init
	logger := log.New(os.Stdout, "[police-service]", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[police-store]", log.LstdFlags)

	// DB init & ping
	store, err := data.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer store.Disconnect(timeoutContext)
	store.Ping()

	// Set LOAD_DB_TEST_DATA to 'false' for persistence between shutdowns
	if os.Getenv("LOAD_DB_TEST_DATA") == "true" {
		err = store.Initialize(context.Background())
		if err != nil {
			logger.Fatalf("Failed to initialize DB: %s", err.Error())
		}
	}

	courtClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	mupClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	ssoClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	court := clients.NewCourtClient(courtClient, os.Getenv("COURT_SERVICE_URI"))
	mup := clients.NewMupClient(mupClient, os.Getenv("MUP_SERVICE_URI"))
	sso := clients.NewSSOClient(ssoClient, os.Getenv("SSO_SERVICE_URI"))

	handler := handlers.NewPoliceHandler(store, court, mup, sso)

	router := mux.NewRouter()
	// Router methods
	router.HandleFunc("/api/v1/traffic-violation/jmbg", handler.GetTrafficViolationsByJMBG).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/traffic-violation", handler.GetAllTrafficViolations).Methods(http.MethodGet)

	authorizedRouter := router.Methods("GET", "POST", "PUT", "DELETE").Subrouter()
	authorizedRouter.HandleFunc("/api/v1/traffic-violation", handler.CreateTrafficViolation).Methods(http.MethodPost)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/{id}", handler.GetTrafficViolationByID).Methods(http.MethodGet)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/{id}", handler.UpdateTrafficViolation).Methods(http.MethodPut)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/{id}", handler.DeleteTrafficViolation).Methods(http.MethodDelete)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/check-all", handler.CheckAll).Methods(http.MethodPost)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/check-alcohol-level", handler.CheckAlcoholLevel).Methods(http.MethodPost)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/check-driver-ban", handler.CheckDriverBan).Methods(http.MethodPost)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/check-driver-permit", handler.CheckDriverPermitValidity).Methods(http.MethodPost)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/check-vehicle-registration", handler.CheckVehicleRegistration).Methods(http.MethodPost)
	authorizedRouter.HandleFunc("/api/v1/traffic-violation/check-vehicle-tire", handler.CheckVehicleTire).Methods(http.MethodPost)

	authorizedRouter.Use(handler.AuthorizeRoles("ADMIN"))

	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	pingRouter := router.Methods("GET").Subrouter()
	pingRouter.HandleFunc("/api/v1", handler.Ping).Methods("GET")
	pingRouter.Use(cors)
	pingRouter.Use(handler.AuthorizeRoles("USER", "ADMIN"))

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
