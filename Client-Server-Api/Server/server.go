package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"context"
	"time"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	  
)

type Cotacao struct {
	Code        string `json:"code"`
	Codein      string `json:"codein"`
	High        string `json:"high"`
	Low         string `json:"low"`
	VarBid      string `json:"varBid"`
	PctChange   string `json:"pctChange"`
	Bid         string `json:"bid"`
	Ask         string `json:"ask"`
	Timestamp   string `json:"timestamp"`
	Create_date string `json:"create_date"`
}

func GetCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	cotacao, error := GetCotacao()
	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(cotacao.Bid)

	

}

func GetCotacao() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	//io.Copy(os.Stdout, res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]Cotacao
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	cotacao, ok := result["USDBRL"]
	if !ok {
		return nil, fmt.Errorf("unexpected JSON structure")
	}
	
    db, err := sql.Open("sqlite3", "./cotacao.db")
    if err != nil {
        return nil, err
    }
    defer db.Close()

	err = InsertCotacao(db, &cotacao)
    if err != nil {
        return nil, err
    }
	
	return &cotacao, nil
}

func InsertCotacao(db *sql.DB, cotacao *Cotacao) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
    defer cancel()

    stmt, err := db.PrepareContext(ctx, "insert into cotacao(code, codein, high, low, varBid, pctChange, bid, ask, timestamp, create_date) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

	_, err = stmt.Exec(cotacao.Code, cotacao.Codein, cotacao.High, cotacao.Low, cotacao.VarBid, cotacao.PctChange, cotacao.Bid, cotacao.Ask, cotacao.Timestamp, cotacao.Create_date)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := sql.Open("sqlite3", "./cotacao.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS cotacao (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "code" TEXT,
        "codein" TEXT,
        "high" TEXT,
        "low" TEXT,
        "varBid" TEXT,
        "pctChange" TEXT,
        "bid" TEXT,
        "ask" TEXT,
        "timestamp" TEXT,
        "create_date" TEXT
    );`
    _, err = db.Exec(createTableSQL)
    if err != nil {
        panic(err)
    }
		http.HandleFunc("/cotacao", GetCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}