package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendSuccessResponse(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Define test data
	type TestResponse struct {
		Message string `json:"message"`
		Status  bool   `json:"status"`
	}
	
	testData := TestResponse{
		Message: "Success",
		Status:  true,
	}
	
	// Call the function to test
	SendSuccessResponse(w, http.StatusOK, testData)
	
	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
	
	// Parse the response body
	var response TestResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error parsing response: %v", err)
	}
	
	// Check response data
	if response.Message != "Success" || !response.Status {
		t.Errorf("Response data doesn't match expected values: %+v", response)
	}
}

func TestSendErrorResponse(t *testing.T) {
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Call the function to test
	errorMessage := "Something went wrong"
	SendErrorResponse(w, http.StatusBadRequest, errorMessage)
	
	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
	
	// Parse the response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Error parsing response: %v", err)
	}
	
	// Check error message
	if response["error"] != errorMessage {
		t.Errorf("Expected error message '%s', got '%s'", errorMessage, response["error"])
	}
}
