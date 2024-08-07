package main

import (
	"context"
	"court/clients"
	"court/data"
	"court/handlers"
	"log"
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
		port = "8083"
	}

	// Context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Logger init
	logger := log.New(os.Stdout, "[court-service]", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[court-store]", log.LstdFlags)

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

	// Client init
	ssoClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	sso := clients.NewSSOClient(ssoClient, os.Getenv("SSO_SERVICE_URI"))

	mupClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	mup := clients.NewMUPClient(mupClient, os.Getenv("MUP_SERVICE_URI"))

	// Handler & router init
	courtHandler := handlers.NewCourtHandler(store, sso, mup)
	router := mux.NewRouter()

	// Router methods
	adminRouter := router.Methods("GET", "POST").Subrouter()
	adminRouter.HandleFunc("/api/v1/get-hearing/{id}", courtHandler.GetCourtHearingByID).Methods("GET")
	adminRouter.HandleFunc("/api/v1/create-hearing-person", courtHandler.CreateHearingPerson).Methods("POST")
	adminRouter.HandleFunc("/api/v1/create-hearing-entity", courtHandler.CreateHearingLegalEntity).Methods("POST")
	adminRouter.HandleFunc("/api/v1/suspensions", courtHandler.CreateSuspension).Methods("POST")
	adminRouter.HandleFunc("/api/v1/warrants", courtHandler.CreateWarrant).Methods("POST")
	adminRouter.HandleFunc("/api/v1/crime-report", courtHandler.RecieveCrimeReport).Methods("POST")
	adminRouter.Use(courtHandler.AuthorizeRoles("ADMIN"))

	authorizedRouter := router.Methods("GET", "PUT").Subrouter()
	authorizedRouter.HandleFunc("/api/v1/courts/{id}", courtHandler.GetCourtByID).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/update-hearing-person", courtHandler.UpdateHearingPerson).Methods("PUT")
	authorizedRouter.HandleFunc("/api/v1/update-hearing-entity", courtHandler.UpdateHearingLegalEntity).Methods("PUT")
	authorizedRouter.HandleFunc("/api/v1/hearings/{jmbg}", courtHandler.GetCourtHearingsByJMBG).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/suspensions/{jmbg}", courtHandler.CheckForSuspension).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/warrants/{jmbg}", courtHandler.CheckForWarrants).Methods("GET")
	authorizedRouter.Use(courtHandler.AuthorizeRoles("USER", "ADMIN"))

	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	pingRouter := router.Methods("GET").Subrouter()
	pingRouter.HandleFunc("/api/v1", courtHandler.Ping).Methods("GET")
	pingRouter.Use(cors)
	pingRouter.Use(courtHandler.AuthorizeRoles("USER", "ADMIN"))

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
