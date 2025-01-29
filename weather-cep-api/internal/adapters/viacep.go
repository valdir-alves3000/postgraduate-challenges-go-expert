package adapters

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"unicode"

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

func GetCityByCEP(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var viacepResp ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viacepResp); err != nil {
		return "", err
	}

	if viacepResp.Erro || viacepResp.City == "" {
		return "", errors.New("zip code not found")
	}

	return removeAccents(viacepResp.City), nil
}
