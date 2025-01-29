package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/handlers"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	http.HandleFunc("/temperature/", handlers.TemperatureHandler)
	http.HandleFunc("/docs", handlers.DocsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs", http.StatusFound)
	})

	log.Printf("Server running on port %s", "8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
