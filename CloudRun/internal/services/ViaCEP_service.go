package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"crypto/tls"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func GetLocationByCEP(cep string) (*ViaCEP, error) {
		client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Get("https://viacep.com.br/ws/" + cep + "/json/")
	
	if err != nil {
		return nil, fmt.Errorf("failed to make request to ViaCEP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response from ViaCEP: %d", resp.StatusCode)
	}

	var viacep ViaCEP
	if err := json.NewDecoder(resp.Body).Decode(&viacep); err != nil {
		return nil, fmt.Errorf("failed to decode response from ViaCEP: %v", err)
	}
	fmt.Printf("Response from ViaCEP: %+v\n", viacep)

	return &viacep, nil
}
