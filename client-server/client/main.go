package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type USDBRL struct {
	Bid string `json:"bid"`
}

func getQuoteFromServer(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server responded with status code %d", resp.StatusCode)
	}

	var result USDBRL
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Bid, nil
}

func saveQuoteToFile(quote string) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", quote))

	return err
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	quote, err := getQuoteFromServer(ctx)
	if err != nil {
		log.Fatalf("Failed to get quote from server: %v", err)
	}

	if err := saveQuoteToFile(quote); err != nil {
		log.Fatalf("Failed to save quote to file: %v", err)
	}

	fmt.Println("Cotação salva em cotacao.txt")
}
