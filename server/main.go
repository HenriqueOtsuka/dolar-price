package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type USDBRL struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createTable()

	http.HandleFunc("/cotacao", searchDolarPriceHandler)
	http.ListenAndServe(":8080", nil)
}

func searchDolarPriceHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	dolarPrice, err := searchDolarPrice(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	saveCtx, saveCancel := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer saveCancel()

	err = insertDolarPrice(saveCtx, dolarPrice.Bid, dolarPrice.Ask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dolarPrice)
}

func searchDolarPrice(ctx context.Context) (*USDBRL, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]USDBRL
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// Access the `USDBRL` key from the map
	dolarPrice, exists := data["USDBRL"]
	if !exists {
		return nil, io.ErrUnexpectedEOF
	}

	return &dolarPrice, nil
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS dolar_price (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT NOT NULL,
		ask TEXT NOT NULL,
		create_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func insertDolarPrice(ctx context.Context, bid, ask string) error {
	query := `INSERT INTO dolar_price (bid, ask) VALUES (?, ?);`
	resultChan := make(chan error, 1)

	go func() {
		_, err := db.Exec(query, bid, ask)
		resultChan <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err() // Timeout or cancellation
	case err := <-resultChan:
		return err
	}
}
