package services

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetTemperature(t *testing.T) {
	originalURL := weatherAPIBaseURL
	defer func() { weatherAPIBaseURL = originalURL }()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{"current": {"temp_c": 25.0, "temp_f": 77.0}}`
		w.Write([]byte(response))
	}))
	defer server.Close()

	weatherAPIBaseURL = server.URL

	os.Setenv("WEATHER_API_KEY", "test-api-key")
	defer os.Unsetenv("WEATHER_API_KEY")

	tempC, tempF, tempK, err := GetTemperature("New York", "NY")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if tempC != 25.0 {
		t.Errorf("Expected tempC to be 25.0, got %f", tempC)
	}
	
	if tempF != 77.0 {
		t.Errorf("Expected tempF to be 77.0, got %f", tempF)
	}
	
	expectedK := 298.0
	if tempK != expectedK {
		t.Errorf("Expected tempK to be %f, got %f", expectedK, tempK)
	}
	
	_, _, _, err = GetTemperature("", "")
	if err == nil {
		t.Error("Expected error with empty parameters, got nil")
	}
}
