package main

import (
	"context"
	"log"
	"mup/data"
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

	router := mux.NewRouter()
	// TODO: Router methods

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
	log.Printf("Server listening on port: %s\n", port)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error while serving request: %v\n", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)
	signal.Notify(sigCh, syscall.SIGTERM)

	sig := <-sigCh
	log.Printf("Recieved terminate, starting gracefull shutdown: %v\n", sig)

	// Gracefull shutdown
	if server.Shutdown(timeoutContext) != nil {
		log.Fatalln("Cannot gracefully shutdown")
	}
	log.Println("Server gracefully stopped")
}
