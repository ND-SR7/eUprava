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

	router.HandleFunc("/api/v1/driving-bans", mupHandler.CheckForPersonsDrivingBans).Methods("GET")
	router.HandleFunc("/api/v1/vehicle", mupHandler.SaveVehicle).Methods("POST")
	router.HandleFunc("/api/v1/registration-request", mupHandler.SubmitRegistrationRequest).Methods("POST")
	router.HandleFunc("/api/v1/approve-registration-request", mupHandler.ApproveRegistration).Methods("POST")
	router.HandleFunc("/api/v1/traffic-permit-request", mupHandler.SubmitTrafficPermitRequest).Methods("POST")
	router.HandleFunc("/api/v1/approve-traffic-permit-request", mupHandler.ApproveTrafficPermitRequest).Methods("POST")

	// For clients
	router.HandleFunc("/api/v1/registered-vehicles", mupHandler.CheckForRegisteredVehicles).Methods("GET")
	router.HandleFunc("/api/v1/driving-ban", mupHandler.IssueDrivingBan).Methods("POST")

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
