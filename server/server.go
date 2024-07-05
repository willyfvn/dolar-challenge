package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/willyfvn/dolar-challenge.git/models"
)

func InitializeServer() string {

	mux := http.NewServeMux()

	mux.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		requestCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		cotacao, err := getCotacao(requestCtx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := models.Cotacao{
			Bid: cotacao,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.ListenAndServe(":8080", mux)

	return "Server initialized"
}

func getCotacao(ctx context.Context) (string, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	var result map[string]map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", nil
	}

	usdbid, ok := result["USDBRL"]["bid"].(string)
	if !ok {
		return "", nil
	}

	return usdbid, nil

}

func InsertCotacao(db *sql.DB, cotacao string) error {
	_, err := db.Exec("INSERT INTO cotacao (bid) VALUES (?)", cotacao)
	if err != nil {
		return err
	}
	fmt.Println("Cotação inserida com sucesso")
	return nil
}
