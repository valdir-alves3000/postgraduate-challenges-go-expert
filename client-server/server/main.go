package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type QuoteResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		Bid        string `json:"bid"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	createTableSQL()
	http.HandleFunc("/cotacao", quoteHandler)
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func quoteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	db, err := sql.Open("sqlite3", "./quotes.db")
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		log.Println("Database connection failed:", err)
		return
	}

	defer db.Close()

	quoteCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	quote, err := getQuote(quoteCtx)
	if err != nil {
		http.Error(w, "Failed to fetch quote", http.StatusInternalServerError)
		log.Println("Failed to fetch quote:", err)
		return
	}

	saveCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	if err := saveQuoteToDB(saveCtx, db, quote); err != nil {
		http.Error(w, "Failed to save quote to database", http.StatusInternalServerError)
		log.Println("Failed to save quote to database:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quote.USDBRL)
}

func getQuote(ctx context.Context) (*QuoteResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var quoteResponse QuoteResponse
	if err := json.NewDecoder(resp.Body).Decode(&quoteResponse); err != nil {
		return nil, err
	}

	return &quoteResponse, nil
}

func saveQuoteToDB(ctx context.Context, db *sql.DB, quote *QuoteResponse) error {
	query := "INSERT INTO quotes (quote, create_date) VALUES (?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, quote.USDBRL.Bid, quote.USDBRL.CreateDate)
	return err
}

func createTableSQL() {
	db, err := sql.Open("sqlite3", "./quotes.db")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS quotes (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"quote" TEXT,
		"create_date" TIMESTAMP
	  );`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal("Failed to create table:", err)
	}
}
