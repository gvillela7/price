package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	repositories "github.com/gvillela7/price/Repositories"
	"github.com/gvillela7/price/config"
	"log"
	"net/http"
	"time"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/cotacao", GetPrice).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func GetPrice(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	priceHost := config.ViperEnvVariable("PRICE_HOST")
	req, err := http.NewRequest("GET", priceHost+"/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req = req.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error: maximum time exceeded -> %v\n", err)
	}
	defer resp.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, bid, err := repositories.Save(resp)
	if err != nil {
		log.Printf("Error: Unable to save to database -> %v\n", err)
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"bid": bid, "statusCode": "ok",
	})
}
