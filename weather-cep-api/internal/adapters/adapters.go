package adapters

import "github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/entities"

type WeatherAdapter interface {
	GetCityByCEP(cep string) (string, error)
	GetTemperature(city string) (entities.TemperatureResponse, error)
}

type DefaultWeatherAdapter struct{}
