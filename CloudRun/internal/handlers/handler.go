package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/Leandroschwab/full-cycle-go/CloudRun/internal/services"
	"github.com/Leandroschwab/full-cycle-go/CloudRun/internal/utils"
)

type CEPCodeRequest struct {
	CEP string `json:"cep"`
}

type CEPCodeResponse struct {
	TemperatureC float64 `json:"temp_C"`
	TemperatureF float64 `json:"temp_F"`
	TemperatureK float64 `json:"temp_K"`
	Error        string  `json:"error,omitempty"`
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
		utils.SendErrorResponse(w, http.StatusNotFound, "Location not found ")
		return
	}

	if location.Localidade == "" || location.Uf == "" {
		fmt.Print("Invalid location data: ", location, "\n")
		utils.SendErrorResponse(w, http.StatusInternalServerError, "can not find zipcode")
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
		TemperatureC: celsius,
		TemperatureF: fahrenheit,
		TemperatureK: kelvinRounded,
	}

	utils.SendSuccessResponse(w, http.StatusOK, response)
}
