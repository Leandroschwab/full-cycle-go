package main

import (
	"log"
	"net/http"

	"github.com/Leandroschwab/full-cycle-go/CloudRun/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/temperature", handlers.HandleCEPCode).Methods("POST")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
