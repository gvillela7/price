package main

import (
	"context"
	"encoding/json"
	"github.com/gvillela7/price/Model"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequest("GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Printf("Error: Wrong URL or server down -> %v\n", err)
	}

	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: maximum time exceeded -> %v\n", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var bid Model.Bid
	if err := json.Unmarshal(body, &bid); err != nil {
		log.Printf("Error: -> %v\n", err)
	}
	text := []byte("Dólar: " + bid.Bid)
	err = os.WriteFile("cotacao.txt", text, 0644)
	if err != nil {
		log.Printf("Error: Could not write to file -> %v\n", err)
	}
}
