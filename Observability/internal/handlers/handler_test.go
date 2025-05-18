package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Leandroschwab/full-cycle-go/Observability/internal/services"
)

type mockLocationService struct{}
type mockTemperatureService struct{}

func (m *mockLocationService) GetLocationByCEP(cep string) (*services.ViaCEP, error) {
	return &services.ViaCEP{
		Localidade: "São Paulo",
		Uf:         "SP",
	}, nil
}

func (m *mockTemperatureService) GetTemperature(city, state string) (float64, float64, float64, error) {
	return 25.0, 77.0, 298.0, nil
}

var (
	originalLocationService    = locationService
	originalTemperatureService = temperatureService
)

func init() {
	locationService = &mockLocationService{}
	temperatureService = &mockTemperatureService{}
}

func TestHandleCEPCode(t *testing.T) {
	defer func() {
		locationService = originalLocationService
		temperatureService = originalTemperatureService
	}()

	requestBody := CEPCodeRequest{CEP: "12345678"}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/temperature", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	HandleCEPCode(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response CEPCodeResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error parsing response: %v", err)
	}

	if response.TemperatureC != 25.0 {
		t.Errorf("Expected temperature in Celsius to be 25.0, got %f", response.TemperatureC)
	}

	if response.TemperatureF != 77.0 {
		t.Errorf("Expected temperature in Fahrenheit to be 77.0, got %f", response.TemperatureF)
	}

	if response.TemperatureK != 298.0 {
		t.Errorf("Expected temperature in Kelvin to be 298.0, got %f", response.TemperatureK)
	}

	invalidRequestBody := CEPCodeRequest{CEP: "123"}
	jsonBody, _ = json.Marshal(invalidRequestBody)

	req = httptest.NewRequest("POST", "/temperature", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()

	HandleCEPCode(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
