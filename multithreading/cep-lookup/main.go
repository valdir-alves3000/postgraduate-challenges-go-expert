package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Message struct {
	id  string
	Msg string
}

type Address struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
}

func handlerAddress(ctx context.Context, url string) (*Address, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var address Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, "", err
	}
	cleanURL := strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://")
	parts := strings.SplitN(cleanURL, "/", 2)
	baseURL := parts[0]

	return &address, baseURL, nil
}

func main() {
	for _, cep := range os.Args[1:] {
		brasilAPIURL := "https://brasilapi.com.br/api/cep/v1/" + cep
		viaCEPURL := "http://viacep.com.br/ws/" + cep + "/json/"
		maxTime := 1.0

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(maxTime)*time.Second)
		defer cancel()

		resultChannelAddress := make(chan struct {
			address *Address
			api     string
			err     error
		})

		go func() {
			address, api, err := handlerAddress(ctx, brasilAPIURL)
			resultChannelAddress <- struct {
				address *Address
				api     string
				err     error
			}{address, api, err}
		}()

		go func() {
			address, api, err := handlerAddress(ctx, viaCEPURL)
			resultChannelAddress <- struct {
				address *Address
				api     string
				err     error
			}{address, api, err}
		}()

		select {
		case res := <-resultChannelAddress:
			if res.err != nil {
				fmt.Println("Erro:", res.err)
				panic(res.err)
			}

			fmt.Printf("  ========================================\n")
			fmt.Printf("Resultado mais rápido vindo de %s:\n", res.api)
			fmt.Printf("CEP: %s\nLogradouro: %s\nBairro: %s\nLocalidade: %s\nUF: %s\n",
				res.address.CEP, res.address.Logradouro, res.address.Bairro, res.address.Localidade, res.address.UF)

		case <-ctx.Done():
			fmt.Printf("Erro: Timeout - Nenhuma das APIs respondeu em %.1f segundo.", maxTime)
		}
	}
}
