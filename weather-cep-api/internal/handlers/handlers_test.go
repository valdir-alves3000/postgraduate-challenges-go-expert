package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/entities"
)

type MockWeatherAdapter struct {
	GetCityByCEPFunc   func(cep string) (string, error)
	GetTemperatureFunc func(city string) (entities.TemperatureResponse, error)
}

func (m *MockWeatherAdapter) GetCityByCEP(cep string) (string, error) {
	if m.GetCityByCEPFunc != nil {
		return m.GetCityByCEPFunc(cep)
	}

	if cep == "01001000" {
		return "São Paulo", nil
	}
	return "", errors.New("zip code not found")
}

func (m *MockWeatherAdapter) GetTemperature(city string) (entities.TemperatureResponse, error) {
	if m.GetTemperatureFunc != nil {
		return m.GetTemperatureFunc(city)
	}

	if city == "São Paulo" {
		return entities.TemperatureResponse{
			City:      "São Paulo",
			Country:   "Brasil",
			Localtime: "2025-01-27 10:00",
			TempC:     25.0,
			TempF:     77.0,
			TempK:     298.0,
		}, nil
	}
	return entities.TemperatureResponse{}, errors.New("temperature not found")
}

func TestTemperatureHandler_Success(t *testing.T) {
	mockAdapter := &MockWeatherAdapter{}
	req := httptest.NewRequest("GET", "/temperature/01001000", nil)
	rr := httptest.NewRecorder()
	handler := TemperatureHandler(mockAdapter)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var response entities.TemperatureResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	expected := entities.TemperatureResponse{
		City:      "São Paulo",
		Country:   "Brasil",
		Localtime: "2025-01-27 10:00",
		TempC:     25.0,
		TempF:     77.0,
		TempK:     298.0,
	}

	if response != expected {
		t.Errorf("Expected response %v, got %v", expected, response)
	}
}

func TestTemperatureHandler_InvalidZipcodeError(t *testing.T) {
	mockAdapter := &MockWeatherAdapter{}
	req := httptest.NewRequest("GET", "/temperature/invalid", nil)
	rr := httptest.NewRecorder()
	handler := TemperatureHandler(mockAdapter)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, status)
	}

	expected := `{"code":422,"message":"invalid zipcode"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
	}
}

func TestTemperatureHandler_CityNotFoundError(t *testing.T) {
	mockAdapter := &MockWeatherAdapter{
		GetCityByCEPFunc: func(cep string) (string, error) {
			return "", errors.New("zip code not found")
		},
	}
	req := httptest.NewRequest("GET", "/temperature/99999999", nil)
	rr := httptest.NewRecorder()
	handler := TemperatureHandler(mockAdapter)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, status)
	}

	expected := `{"code":404,"message":"can not find zipcode"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
	}
}

func TestTemperatureHandler_TemperatureNotFoundError(t *testing.T) {
	mockAdapter := &MockWeatherAdapter{
		GetCityByCEPFunc: func(cep string) (string, error) {
			return "São Paulo", nil
		},
		GetTemperatureFunc: func(city string) (entities.TemperatureResponse, error) {
			return entities.TemperatureResponse{}, errors.New("temperature not found")
		},
	}
	req := httptest.NewRequest("GET", "/temperature/01001000", nil)
	rr := httptest.NewRecorder()
	handler := TemperatureHandler(mockAdapter)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	expected := `{"code":500,"message":"can not find temperature"}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("Expected body %v, got %v", expected, rr.Body.String())
	}
}
