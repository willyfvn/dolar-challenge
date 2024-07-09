package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/willyfvn/dolar-challenge.git/db"
	"github.com/willyfvn/dolar-challenge.git/models"
)

func InitializeServer() string {

	mux := http.NewServeMux()

	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, req *http.Request) {
		mydb := db.StartDb()
		defer mydb.Close()

		requestCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		dbCtx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
		defer cancel()

		cotacao, err := getCotacao(requestCtx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusRequestTimeout)
			fmt.Println(err)
			return

		}
		err = InsertCotacao(mydb, cotacao, dbCtx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusRequestTimeout)
			return
		}

		response := models.Cotacao{
			Bid: cotacao,
		}
		if response.Bid == "" {
			http.Error(w, "bid not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":8080", mux)

	return "Server initialized"
}

func getCotacao(ctx context.Context) (string, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}
	usdbid, ok := result["USDBRL"]["bid"].(string)
	if !ok {
		return "", fmt.Errorf("bid not found")
	}
	fmt.Println(usdbid)

	return "usdbid", nil

}

func InsertCotacao(db *sql.DB, cotacao string, ctx context.Context) error {
	_, err := db.ExecContext(ctx, "INSERT INTO cotacao (bid) VALUES (?)", cotacao)

	if err != nil {
		log.Println("error inserting cotacao:")
		log.Println(err)
		return err
	}
	fmt.Println("Cotação inserida com sucesso")
	return nil
}
