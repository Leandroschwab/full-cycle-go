package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Leandroschwab/full-cycle-go/Observability/internal/services"
	"github.com/Leandroschwab/full-cycle-go/Observability/internal/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
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

type Handler struct {
	TemplateData *TemplateData
}

type TemplateData struct {
	Funtion    string
	HTTP_PORT  string
	OTELTracer trace.Tracer
}

func NewServer(templateData *TemplateData) *Handler {
	return &Handler{
		TemplateData: templateData,
	}
}

var tracer = otel.Tracer("handlers")

func HandleCEPCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, span := tracer.Start(ctx, "HandleCEPCode")
	defer span.End()

	var request CEPCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.RecordError(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cep := strings.TrimSpace(request.CEP)
	if len(cep) != 8 {
		span.AddEvent("Invalid CEP length")
		utils.SendErrorResponse(w, http.StatusBadRequest, "invalid zipcode")
		return
	}
	location, err := locationService.GetLocationByCEP(cep)
	if err != nil {
		span.RecordError(err)
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
		span.RecordError(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch temperature")
		return
	}
	kelvinRounded := math.Round(kelvin*10) / 10
	response := CEPCodeResponse{
		City:         location.Localidade,
		TemperatureC: celsius,
		TemperatureF: fahrenheit,
		TemperatureK: kelvinRounded,
	}
	utils.SendSuccessResponse(w, http.StatusOK, response)
}

func (h *Handler) ValidateCEPCode(w http.ResponseWriter, r *http.Request) {
	// Extract context from the incoming request.
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := otel.GetTextMapPropagator().Extract(r.Context(), carrier)

	// Start a span for the validation operation.
	ctx, span := h.TemplateData.OTELTracer.Start(ctx, "ValidateCEPCode")
	defer span.End()

	var request CEPCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.RecordError(err)
		utils.SendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cep := strings.TrimSpace(request.CEP)
	if len(cep) != 8 {
		span.AddEvent("Invalid CEP length")
		utils.SendErrorResponse(w, http.StatusUnprocessableEntity, "invalid zipcode")
		return
	}

	// Start a child span for the communication with Service B (orchestrator)
	ctx, callSpan := h.TemplateData.OTELTracer.Start(ctx, "CallOrchestrator")
	defer callSpan.End()

	httpClient := &http.Client{}
	orchestrator_URL := os.Getenv("ORCHSTRATOR_URL")
	orchestrator_PORT := os.Getenv("ORCHSTRATOR_PORT")
	url := fmt.Sprintf("http://%s:%s/temperature", orchestrator_URL, orchestrator_PORT)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(`{"cep":"`+cep+`"}`))
	if err != nil {
		callSpan.RecordError(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create request to orchestrator")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := httpClient.Do(req)
	if err != nil {
		callSpan.RecordError(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to connect to orchestrator service")
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			callSpan.RecordError(err)
			utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to parse orchestrator error response")
			return
		}
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	var orchestratorResponse CEPCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&orchestratorResponse); err != nil {
		callSpan.RecordError(err)
		utils.SendErrorResponse(w, http.StatusInternalServerError, "Failed to parse orchestrator response")
		return
	}

	//add 500ms delay

	time.Sleep(500 * time.Millisecond)
	 
	w.WriteHeader(resp.StatusCode)
	json.NewEncoder(w).Encode(orchestratorResponse)
}
