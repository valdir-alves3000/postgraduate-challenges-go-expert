package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"unicode"

	"go.opentelemetry.io/otel"
	"golang.org/x/text/unicode/norm"
)

type ViaCEPResponse struct {
	City string `json:"localidade"`
	Erro bool   `json:"erro,omitempty"`
}

func removeAccents(input string) string {
	t := norm.NFD.String(input)
	var result []rune
	for _, r := range t {
		if !unicode.Is(unicode.Mn, r) {
			result = append(result, r)
		}
	}
	return string(result)
}

func (d *DefaultWeatherAdapter) GetCityByCEP(cep string) (string, error) {
	ctx := context.Background()
	tracer := otel.Tracer("weather-cep-api")
	_, span := tracer.Start(ctx, "GetCityByCEP")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)

	resp, err := http.Get(url)
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	defer resp.Body.Close()

	var viacepResp ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viacepResp); err != nil {
		span.RecordError(err)
		return "", err
	}

	if viacepResp.Erro || viacepResp.City == "" {
		err := errors.New("zip code not found")
		span.RecordError(err)
		return "", err
	}

	return removeAccents(viacepResp.City), nil
}
