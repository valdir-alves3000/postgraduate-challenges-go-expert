package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/entities"
	"go.opentelemetry.io/otel"
)

type WeatherAPIResponse struct {
	Location struct {
		Name      string `json:"name"`
		Country   string `json:"country"`
		Localtime string `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (d *DefaultWeatherAdapter) GetTemperature(city string) (entities.TemperatureResponse, error) {
	ctx := context.Background()
	tracer := otel.Tracer("weather-cep-api")
	_, span := tracer.Start(ctx, "GetTemperature")
	defer span.End()

	apiKey := os.Getenv("WEATHER_API_KEY")
	city = strings.ReplaceAll(city, " ", "%20")

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, city)

	resp, err := http.Get(url)
	if err != nil {
		span.RecordError(err)
		log.Printf("Error fetching temperature from Weather API: %v", err)
		return entities.TemperatureResponse{}, err
	}
	defer resp.Body.Close()

	var weatherResp WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		span.RecordError(err)
		log.Printf("Error decoding JSON response from Weather API: %v", err)
		return entities.TemperatureResponse{}, err
	}

	return entities.TemperatureResponse{
		City:      weatherResp.Location.Name,
		Country:   weatherResp.Location.Country,
		Localtime: weatherResp.Location.Localtime,
		TempC:     weatherResp.Current.TempC,
		TempF:     math.Round((weatherResp.Current.TempC*1.8+32)*10) / 10,
		TempK:     math.Round((weatherResp.Current.TempC+273)*10) / 10,
	}, nil
}
