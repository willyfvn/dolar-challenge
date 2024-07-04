package client

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/willyfvn/dolar-challenge.git/db"
	"github.com/willyfvn/dolar-challenge.git/models"
)

func FetchCotacao() string {

	mydb := db.StartDb()
	defer mydb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return err.Error()
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}

	defer resp.Body.Close()

	cotacao := models.Cotacao{}
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		return err.Error()
	}

	db.InsertCotacao(mydb, cotacao.Bid)

	err = saveCotacao(&cotacao)
	if err != nil {
		return err.Error()
	}

	return "Cotação salva com sucesso"

}

func saveCotacao(cotacao *models.Cotacao) error {
	content := "Dólar: " + cotacao.Bid
	return os.WriteFile("cotacao.txt", []byte(content), 0644)
}
