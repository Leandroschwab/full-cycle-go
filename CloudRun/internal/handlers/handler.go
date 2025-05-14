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

func HandleCEPCode(w http.ResponseWriter, r *http.Request) {
	var request CEPCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cep := strings.TrimSpace(request.CEP)
	if len(cep) != 8 {
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid CEP code format")
		return
	}

	location, err := services.GetLocationByCEP(cep)
	if err != nil {
		fmt.Print("Error fetching location: ", err, "\n")
		utils.SendErrorResponse(w, http.StatusNotFound, "Location not found ")
		return
	}

	celsius, fahrenheit, kelvin, err := services.GetTemperature(location.Localidade, location.Uf)
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
