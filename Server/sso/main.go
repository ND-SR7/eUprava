package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sso/data"
	"sso/handlers"
	"syscall"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	// Context
	timeoutContext, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Logger init
	logger := log.New(os.Stdout, "[sso-service]", log.LstdFlags)
	storeLogger := log.New(os.Stdout, "[sso-store]", log.LstdFlags)

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

	// Handler & router init
	ssoHandler := handlers.NewSSOHandler(store)
	router := mux.NewRouter()

	// Router methods
	router.HandleFunc("/api/v1/login", ssoHandler.Login).Methods("POST")
	router.HandleFunc("/api/v1/register-person", ssoHandler.RegisterPerson).Methods("POST")
	router.HandleFunc("/api/v1/register-entity", ssoHandler.RegisterLegalEntity).Methods("POST")
	router.HandleFunc("/api/v1/activate/{activationCode}", ssoHandler.ActivateAccount).Methods("GET")
	router.HandleFunc("/api/v1/recover-password", ssoHandler.RecoverPassword).Methods("POST")
	router.HandleFunc("/api/v1/reset-password", ssoHandler.ResetPassword).Methods("POST")

	authorizedRouter := router.Methods("GET").Subrouter()
	authorizedRouter.HandleFunc("/api/v1/user/{accountID}", ssoHandler.GetUserByAccountID).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/user/email/{email}", ssoHandler.GetUserByEmail).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/user/jmbg/{jmbg}", ssoHandler.GetPersonByJMBG).Methods("GET")
	authorizedRouter.HandleFunc("/api/v1/user/mb/{mb}", ssoHandler.GetLegalEntityByMB).Methods("GET")
	authorizedRouter.Use(ssoHandler.AuthorizeRoles("USER", "ADMIN"))

	cors := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"*"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
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
