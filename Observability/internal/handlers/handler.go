package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/Leandroschwab/full-cycle-go/Observability/internal/services"
	"github.com/Leandroschwab/full-cycle-go/Observability/internal/utils"
)

type CEPCodeRequest struct {
	CEP string `json:"cep"`
}

type CEPCodeResponse struct {
	City         string  `json:"city"`
	TemperatureC float64 `json:"temp_C"`
	TemperatureF float64 `json:"temp_F"`
	TemperatureK float64 `json:"temp_K"`
	Error        string  `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LocationService interface {
	GetLocationByCEP(cep string) (*services.ViaCEP, error)
}

type TemperatureService interface {
	GetTemperature(city, state string) (float64, float64, float64, error)
}

type defaultLocationService struct{}
type defaultTemperatureService struct{}

func (s *defaultLocationService) GetLocationByCEP(cep string) (*services.ViaCEP, error) {
	return services.GetLocationByCEP(cep)
}

func (s *defaultTemperatureService) GetTemperature(city, state string) (float64, float64, float64, error) {
	return services.GetTemperature(city, state)
}

var (
	locationService    LocationService    = &defaultLocationService{}
	temperatureService TemperatureService = &defaultTemperatureService{}
)

func GetLocationService() LocationService {
	return locationService
}

func SetLocationService(service LocationService) {
	locationService = service
}

func GetTemperatureService() TemperatureService {
	return temperatureService
}

func SetTemperatureService(service TemperatureService) {
	temperatureService = service
}

func HandleCEPCode(w http.ResponseWriter, r *http.Request) {
	var request CEPCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cep := strings.TrimSpace(request.CEP)
	if len(cep) != 8 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid zipcode")
		return
	}

	location, err := locationService.GetLocationByCEP(cep)
	if err != nil {
		fmt.Print("Error fetching location: ", err, "\n")
		utils.SendErrorResponse(w, http.StatusNotFound, "can not find zipcode")
		return
	}

	if location.Localidade == "" || location.Uf == "" {
		fmt.Print("Invalid location data: ", location, "\n")
		utils.SendErrorResponse(w, http.StatusNotFound, "can not find zipcode")
		return
	}

	celsius, fahrenheit, kelvin, err := temperatureService.GetTemperature(location.Localidade, location.Uf)
	if err != nil {
		fmt.Print("Error fetching temperature: ", err, "\n")
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch temperature")
		return
	}

	// Ensure the Kelvin value is properly rounded to 1 decimal place
	kelvinRounded := math.Round(kelvin*10) / 10

	response := CEPCodeResponse{
		City:         location.Localidade,
		TemperatureC: celsius,
		TemperatureF: fahrenheit,
		TemperatureK: kelvinRounded,
	}

	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func ValidateCEPCode(w http.ResponseWriter, r *http.Request) {
	var request CEPCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cep := strings.TrimSpace(request.CEP)
	if len(cep) != 8 {
		fmt.Print("Invalid CEP code: ", cep, "\n")
		utils.SendErrorResponse(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}
	fmt.Print("CEP Valid: ", cep, "\n")
	//CEP Valido send request to serviceb
	httpClient := &http.Client{}
	orchestrator_URL := os.Getenv("ORCHSTRATOR_URL")
	orchestrator_PORT := os.Getenv("ORCHSTRATOR_PORT")
	url := fmt.Sprintf("http://%s:%s/temperature", orchestrator_URL, orchestrator_PORT)

	req, err := http.NewRequest("POST", url, strings.NewReader(`{"cep":"`+cep+`"}`))
	if err != nil {
		fmt.Println("Error creating request:", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create request to orchestrator")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error making request to orchestrator:", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to connect to orchestrator service")
		return
	}
	defer resp.Body.Close()

	// Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// For error responses, we should only include the error message
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			fmt.Println("Error decoding orchestrator error response:", err)
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to parse orchestrator error response")
			return
		}

		// Return only the error message with the same status code
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(errorResp)
		fmt.Println("Error retrieving CEP code from orchestrator")
		return
	}

	// For successful responses, return the temperature data
	var orchestratorResponse CEPCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&orchestratorResponse); err != nil {
		fmt.Println("Error decoding orchestrator success response:", err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to parse orchestrator response")
		return
	}

	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(orchestratorResponse)
	fmt.Println("Successfully retrieved CEP code from orchestrator")
}
