package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type USDBRL struct {
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	CreateDate string `json:"create_date"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Request timed out")
		} else {
			fmt.Println("Error sending request:", err)
		}
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received non-200 response: %s\n", res.Status)
		return
	}

	var usdbrl USDBRL
	err = json.NewDecoder(res.Body).Decode(&usdbrl)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("./cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Dolar: %s\n", usdbrl.Bid)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Dolar price has been written")
}
