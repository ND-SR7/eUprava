package main

import (
	"context"
	"log"
	"mup/clients"
	"mup/data"
	"mup/handlers"
	"mup/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8081"
	}

	// Context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Logger init
	logger := log.New(os.Stdout, "[mup-service] ", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[mup-store] ", log.LstdFlags)

	// DB init & ping
	store, err := data.New(timeoutContext, storeLogger)
	if err != nil {
		logger.Fatal(err)
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

	ssoClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	court := clients.NewCourtClient(courtClient, os.Getenv("COURT_SERVICE_URI"))
	sso := clients.NewSSOClient(ssoClient, os.Getenv("SSO_SERVICE_URI"))

	mupService := services.NewMupService(store, storeLogger, sso, court)
	mupHandler := handlers.NewMupHandler(mupService, storeLogger)

	router := mux.NewRouter()

	//GET
	router.HandleFunc("/api/v1/persons-vehicles", mupHandler.GetPersonsVehicles).Methods("GET")
	router.HandleFunc("/api/v1/driving-bans", mupHandler.CheckForPersonsDrivingBans).Methods("GET")
	router.HandleFunc("/api/v1/persons-registrations", mupHandler.GetPersonsRegistrations).Methods("GET")
	router.HandleFunc("/api/v1/persons-driving-permit", mupHandler.GetUserDrivingPermit).Methods("GET")

	//POST
	router.HandleFunc("/api/v1/vehicle", mupHandler.SaveVehicle).Methods("POST")
	router.HandleFunc("/api/v1/registration-request", mupHandler.SubmitRegistrationRequest).Methods("POST")
	router.HandleFunc("/api/v1/traffic-permit-request", mupHandler.SubmitTrafficPermitRequest).Methods("POST")

	authorizedRouter := router.Methods("GET", "POST").Subrouter()
	authorizedRouter.HandleFunc("/api/v1/pending-registration-requests", mupHandler.GetPendingRegistrationRequests).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/pending-traffic-permit-requests", mupHandler.GetPendingTrafficPermitRequests).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/approve-registration-request", mupHandler.ApproveRegistration).Methods("POST")
	authorizedRouter.HandleFunc("/api/v1/approve-traffic-permit-request", mupHandler.ApproveTrafficPermitRequest).Methods("POST")

	// For clients
	authorizedRouter.HandleFunc("/api/v1/registered-vehicles", mupHandler.CheckForRegisteredVehicles).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/driving-ban", mupHandler.IssueDrivingBan).Methods("POST")
	authorizedRouter.HandleFunc("/api/v1/registration-by-plate", mupHandler.GetRegistrationByPlate).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/check-persons-driving-ban", mupHandler.GetDrivingBan).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/check-persons-driving-permit", mupHandler.GetDrivingPermitByJMBG).Methods("GET")
	authorizedRouter.Use(mupHandler.AuthorizeRoles("ADMIN"))

	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	pingRouter := router.Methods("GET").Subrouter()
	pingRouter.HandleFunc("/api/v1", mupHandler.Ping).Methods("GET")
	pingRouter.Use(cors)
	pingRouter.Use(mupHandler.AuthorizeRoles("USER", "ADMIN"))

	mupService.SaveMup()

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
