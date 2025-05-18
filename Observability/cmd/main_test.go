package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Leandroschwab/full-cycle-go/Observability/internal/handlers"
	"github.com/Leandroschwab/full-cycle-go/Observability/internal/services"
	"github.com/gorilla/mux"
)

// Mock implementations for services
type mockLocationService struct{}
type mockTemperatureService struct{}

func (m *mockLocationService) GetLocationByCEP(cep string) (*services.ViaCEP, error) {
	return &services.ViaCEP{
		Localidade: "Rio de Janeiro",
		Uf:         "RJ",
	}, nil
}

func (m *mockTemperatureService) GetTemperature(city, state string) (float64, float64, float64, error) {
	return 25.0, 77.0, 298.2, nil
}

func setupTestRouter() *mux.Router {
	// Save original services
	originalLocationService := handlers.GetLocationService()
	originalTemperatureService := handlers.GetTemperatureService()

	// Override with mock services
	handlers.SetLocationService(&mockLocationService{})
	handlers.SetTemperatureService(&mockTemperatureService{})

	// Create router with the temperature endpoint
	router := mux.NewRouter()
	router.HandleFunc("/temperature", handlers.HandleCEPCode).Methods("POST")

	// Restore original services when test is done
	// We're not actually restoring here because the test is short-lived,
	// but in real code you might want to use defer
	_ = originalLocationService
	_ = originalTemperatureService

	return router
}

func TestTemperatureEndpoint(t *testing.T) {
	router := setupTestRouter()

	requestBody := handlers.CEPCodeRequest{CEP: "12345678"}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/temperature", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected HTTP 200 response, got %d", rr.Code)
	}

	req = httptest.NewRequest("GET", "/temperature", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected HTTP 405 Method Not Allowed for GET request, got %d", rr.Code)
	}

	req = httptest.NewRequest("POST", "/invalid-endpoint", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected HTTP 404 Not Found for invalid endpoint, got %d", rr.Code)
	}
}
