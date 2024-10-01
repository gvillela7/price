package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gvillela7/price/Model"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
	"time"
)

var db *sql.DB

func init() {
	database, err := sql.Open("sqlite3", "Database/price.db")
	if err != nil {
		log.Fatal(err)
	}
	db = database
}

func Save(resp *http.Response) (int64, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	stmt, err := db.Prepare(
		`INSERT INTO prices
    				(code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)
				VALUES 
				    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
				`,
	)
	if err != nil {
		return int64(0), "", err
	}

	body, _ := io.ReadAll(resp.Body)
	var price Model.Price
	if err := json.Unmarshal(body, &price); err != nil {
		return int64(0), "", err
	}

	res, err := stmt.ExecContext(ctx,
		price.USDBRL.Code, price.USDBRL.Codein, price.USDBRL.Name, price.USDBRL.High, price.USDBRL.Low,
		price.USDBRL.VarBid, price.USDBRL.PctChange, price.USDBRL.Bid, price.USDBRL.Ask, price.USDBRL.Timestamp,
		price.USDBRL.CreateDate,
	)
	if err != nil {
		return int64(0), "", err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return int64(0), "", err
	}
	return lastID, fmt.Sprintf("%s", price.USDBRL.Bid), nil
}
