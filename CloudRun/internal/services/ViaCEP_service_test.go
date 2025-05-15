package services

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var originalHttpClient = func() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func TestGetLocationByCEP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		response := `{
			"cep": "22021-001",
			"logradouro": "Avenida Atlântica",
			"complemento": "de 1662 a 2172 - lado par",
			"unidade": "",
			"bairro": "Copacabana",
			"localidade": "Rio de Janeiro",
			"uf": "RJ",
			"estado": "Rio de Janeiro",
			"regiao": "Sudeste",
			"ibge": "3304557",
			"gia": "",
			"ddd": "21",
			"siafi": "6001"
		}`
		
		w.Write([]byte(response))
	}))
	defer server.Close()

	originalURL := "https://viacep.com.br/ws/"
	
	httpClientCreator = func() *http.Client {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
	
	viacepBaseURL = server.URL + "/ws/"
	
	defer func() {
		httpClientCreator = originalHttpClient
		viacepBaseURL = originalURL
	}()

	location, err := GetLocationByCEP("22021001")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if location.Localidade != "Rio de Janeiro" {
		t.Errorf("Expected city to be Rio de Janeiro, got %s", location.Localidade)
	}
	
	if location.Uf != "RJ" {
		t.Errorf("Expected state to be RJ, got %s", location.Uf)
	}
	
	if location.Bairro != "Copacabana" {
		t.Errorf("Expected neighborhood to be Copacabana, got %s", location.Bairro)
	}
}
