package main

import (
	"net/http"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/APIs/configs"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/APIs/internal/entity"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/APIs/internal/infra/database"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/APIs/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	http.HandleFunc("/products", productHandler.CreateProduct)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.ListenAndServe(":8000", nil)
}
