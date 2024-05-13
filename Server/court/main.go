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

	// Client init
	ssoClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     10,
		},
	}

	sso := clients.NewSSOClient(ssoClient, os.Getenv("SSO_SERVICE_URI"))

	// Handler & router init
	courtHandler := handlers.NewCourtHandler(store, sso)
	router := mux.NewRouter()

	// Router methods
	router.HandleFunc("/api/v1/get-hearing/{id}", courtHandler.GetCourtHearingByID).Methods("GET")
	router.HandleFunc("/api/v1/create-hearing-person", courtHandler.CreateHearingPerson).Methods("POST")
	router.HandleFunc("/api/v1/create-hearing-entity", courtHandler.CreateHearingLegalEntity).Methods("POST")
	router.HandleFunc("/api/v1/update-hearing-person", courtHandler.UpdateHearingPerson).Methods("PUT")
	router.HandleFunc("/api/v1/update-hearing-entity", courtHandler.UpdateHearingLegalEntity).Methods("PUT")

	authorizedRouter := router.Methods("POST").Subrouter()
	authorizedRouter.HandleFunc("/api/v1/user/crime-report", courtHandler.RecieveCrimeReport).Methods("POST")
	authorizedRouter.Use(courtHandler.AuthorizeRoles("ADMIN"))

	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type"}),
	)

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
