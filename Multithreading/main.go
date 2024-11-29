package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Message struct {
	Endereco string
}

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEP struct {
	Cep         string `json: cep`
	Logradouro  string `json: logradouro`
	Complemento string `json: complemento`
	Unidade     string `json: unidade`
	Bairro      string `json: bairro`
	Localidade  string `json: localidade`
	Uf          string `json: uf`
	Estado      string `json: estado`
	Regiao      string `json: regiao`
	Ibge        string `json: ibge`
	Gia         string `json: gia`
	Ddd         string `json: ddd`
	Siafi       string `json: siafi`
}

func main() {
	c1 := make(chan Message)
	c2 := make(chan Message)

	go reqViaCEP(c2, "01153000 ")
	go reqBrasilAPI(c1, "01153000 ")

	for {
		select {
		case msg := <-c1: // rabbitmq
			fmt.Printf("Received from BrasilAPI: Endereco: %s\n", msg.Endereco)
			return

		case msg := <-c2: // kafka
			fmt.Printf("Received from ViaCEP: Endereco: %s\n", msg.Endereco)
			return

		case <-time.After(time.Second * 1):
			println("timeout")
			return
		}
	}
}

func reqBrasilAPI(c1 chan Message, cep string) {
	//time.Sleep(500 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	var result BrasilAPI
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("BrasilAPI " + result.Street + ", " + result.Neighborhood + ", " + result.City + ", " + result.State)
	c1 <- Message{result.Street + ", " + result.Neighborhood + ", " + result.City + ", " + result.State}
}

func reqViaCEP(c2 chan Message, cep string) {
	//time.Sleep(500 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	var result ViaCEP
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("ViaCEP " + result.Logradouro + ", " + result.Bairro + ", " + result.Localidade + ", " + result.Uf)
	c2 <- Message{result.Logradouro + ", " + result.Bairro + ", " + result.Localidade + ", " + result.Uf}

}
