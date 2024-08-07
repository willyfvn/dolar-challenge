package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/willyfvn/dolar-challenge.git/models"
)

func FetchCotacao() string {

	requestCtx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(requestCtx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return err.Error()
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Sprintf("Error creating request: %s", resp.Status)
	}

	cotacao := models.Cotacao{}
	err = json.NewDecoder(resp.Body).Decode(&cotacao)
	if err != nil {
		return err.Error()
	}

	err = saveCotacao(&cotacao)
	if err != nil {
		return fmt.Sprintf("4 Error creating request: %s", err)
	}

	return "Cotação salva com sucesso"

}

func saveCotacao(cotacao *models.Cotacao) error {
	content := "Dólar: " + cotacao.Bid
	return os.WriteFile("cotacao.txt", []byte(content), 0644)
}
