package main

import (
	"log"
	"net/http"

	"os"

	"github.com/Leandroschwab/full-cycle-go/Observability/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	service := os.Getenv("FUNCTION")
	switch service {
	case "orchestrator":
		router.HandleFunc("/temperature", handlers.HandleCEPCode).Methods("POST")
	case "inputvalidator":
		router.HandleFunc("/temperature", handlers.ValidateCEPCode).Methods("POST")
	default:
		log.Fatal("Invalid function specified")
	}
	port := os.Getenv("HTTP_PORT")

	log.Println("Starting server on ", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)

	}
}
