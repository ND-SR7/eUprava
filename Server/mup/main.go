package main

import (
	"context"
	"log"
	"mup/data"
	"mup/handlers"
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

	mupHandler := handlers.NewMupHandler(store, storeLogger)

	router := mux.NewRouter()

	checkPersonsDrivingBans := router.Methods(http.MethodGet).Path("/drivingBans").Subrouter()
	checkPersonsDrivingBans.HandleFunc("", mupHandler.CheckForPersonsDrivingBans)

	saveVehicle := router.Methods(http.MethodPost).Path("/vehicle").Subrouter()
	saveVehicle.HandleFunc("", mupHandler.SaveVehicle)
	//createPersonRouter.Use(mupHandler.MiddlewarePersonDeserialization)

	issueDrivingBan := router.Methods(http.MethodPost).Path("/drivingBan").Subrouter()
	issueDrivingBan.HandleFunc("", mupHandler.IssueDrivingBan)

	submitRegistrationRequest := router.Methods(http.MethodPost).Path("/registrationRequest").Subrouter()
	submitRegistrationRequest.HandleFunc("", mupHandler.SubmitRegistrationRequest)

	approveRegistrationRequest := router.Methods(http.MethodPost).Path("/approveRegistrationRequest").Subrouter()
	approveRegistrationRequest.HandleFunc("", mupHandler.ApproveRegistration)

	submitTrafficPermitRequest := router.Methods(http.MethodPost).Path("/trafficPermitRequest").Subrouter()
	submitTrafficPermitRequest.HandleFunc("", mupHandler.SubmitTrafficPermitRequest)

	approveTrafficPermitRequest := router.Methods(http.MethodPost).Path("/approveTrafficPermitRequest").Subrouter()
	approveTrafficPermitRequest.HandleFunc("", mupHandler.ApproveTrafficPermitRequest)

	////Save mup
	//mupID, err := primitive.ObjectIDFromHex("607d22b837ede6b71eef3e82")
	//if err == nil {
	//	address := data.Address{
	//		Municipality: "ss",
	//		Locality:     "Novi Sad",
	//		StreetName:   "Dunavska",
	//		StreetNumber: 1,
	//	}
	//	mup := data.Mup{
	//		ID:             mupID,
	//		Name:           "Mup",
	//		Address:        address,
	//		Vehicles:       make([]primitive.ObjectID, 0),
	//		TrafficPermits: make([]primitive.ObjectID, 0),
	//		Plates:         make([]string, 0),
	//		DrivingBans:    make([]primitive.ObjectID, 0),
	//		Registrations:  make([]string, 0),
	//	}
	//	err = mupHandler.SaveMup(mup)
	//	if err != nil {
	//		log.Printf("Failed to save mup: %v", err)
	//	}
	//	if err == nil {
	//		log.Printf("Saved mup: %v", mup)
	//	}
	//}

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
