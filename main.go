package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lemm8/IDS-auth-service/cfg"
)

func init() {
	cfg.LoadEnvVariables()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(w, "Hello World")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(w, "Auth Service")
}

func main() {
	// Logger
	l := log.New(os.Stdout, "auth-service-logger", log.LstdFlags)

	// Create ServeMux
	serveMux := mux.NewRouter()

	// Register Handlers
	getRouter := serveMux.Methods("GET").Subrouter()
	getRouter.HandleFunc("/", homeHandler)
	getRouter.HandleFunc("/auth", authHandler)

	// CORS
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	serverAddr := "127.0.0.1:" + os.Getenv("SERVER_PORT")

	// Create custom server
	server := http.Server{
		Addr:         serverAddr,
		Handler:      corsHandler(serveMux),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	// Handle ListenAndServe in goroutine to avoid blocking
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal("Error: ", err)
		}
	}()

	// Broadcast message when interrupt or kill happens
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	l.Println("Received terminate, graceful shutdown", sig)

	// Run server
	server.ListenAndServe()

	// Graceful shutdown
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}
