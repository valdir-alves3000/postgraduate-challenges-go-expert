package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/adapters"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/internal_error"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/weather-cep-api/internal/validations"
)

func TemperatureHandler(adapter adapters.WeatherAdapter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		zipcode := r.URL.Path[len("/temperature/"):]

		isZipcodeValid, formattedZipcode := validations.IsZipcodeValid(zipcode)

		if !isZipcodeValid {
			log.Printf("%s: %v", "Invalid ZIP code", zipcode)
			internal_error.InvalidZipcodeError(w)
			return
		}

		city, err := adapter.GetCityByCEP(formattedZipcode)
		if err != nil {
			log.Printf("%s: %v", "Failed to find city for ZIP code", err)
			internal_error.CityNotFoundError(w)
			return
		}

		temp, err := adapter.GetTemperature(city)
		if err != nil {
			log.Printf("%s: %v", "Failed to retrieve temperature", err)
			internal_error.TemperatureNotFoundError(w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(temp)
	}
}

func DocsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
<!DOCTYPE html>
<html lang="pt-BR">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Documentação da API - Weather CEP</title>
    <style>
      :root {
        --primary-color: #2563eb;
        --text-color: #374151;
        --bg-code: #f3f4f6;
        --border-color: #e5e7eb;
      }

      body {
        font-family: "Segoe UI", sans-serif;
        line-height: 1.6;
        color: var(--text-color);
        max-width: 800px;
        margin: 0 auto;
        padding: 2rem;
        background-color: #fff;
      }

      h1 {
        color: var(--primary-color);
        font-size: 2.25rem;
        margin-bottom: 1.5rem;
        border-bottom: 2px solid var(--border-color);
        padding-bottom: 0.5rem;
      }

      h2 {
        color: var(--text-color);
        font-size: 1.5rem;
        margin: 2rem 0 1rem;
      }

      .intro {
        background-color: #f9fafb;
        border-radius: 8px;
        padding: 1.5rem;
        margin-bottom: 2rem;
        border-left: 4px solid var(--primary-color);
      }

      .endpoint {
        background-color: white;
        border: 1px solid var(--border-color);
        border-radius: 8px;
        padding: 1.5rem;
        margin-bottom: 1.5rem;
      }

      code {
        background-color: var(--bg-code);
        padding: 0.2rem 0.4rem;
        border-radius: 4px;
        font-family: "Courier New", Courier, monospace;
        font-size: 0.9rem;
      }

      .endpoint-title {
        font-weight: bold;
        color: var(--primary-color);
        margin-bottom: 0.5rem;
      }

      .params {
        margin-top: 1rem;
      }

      .param {
        margin-bottom: 0.5rem;
      }

      .param-name {
        font-weight: bold;
      }

      .example {
        background-color: var(--bg-code);
        padding: 4px;
        border-radius: 4px;
        margin-top: 0.5rem;
      }
    </style>
  </head>
  <body>
    <h1>Weather CEP API</h1>

    <div class="intro">
      <p>Bem-vindo à Weather CEP API!</p>
      <p>
        Esta documentação fornece informações sobre como utilizar nossa API para
        obter a previsão do tempo e outras informações meteorológicas com base em CEPs.
      </p>
    </div>

    <h2>Endpoints Disponíveis</h2>

    <div class="endpoint">
      <div class="endpoint-title">Obter Previsão do Tempo</div>
      <code>GET /temperature/{cep}</code>
      <p>
        Retorna informações sobre a previsão do tempo para o CEP fornecido.
      </p>

      <div class="params">
        <div class="param">
          <span class="param-name">cep</span>: CEP para o qual deseja obter a previsão.
        </div>
        <div class="example">
          Exemplo:
          <code>/temperature/01001000</code>
        </div>
      </div>
    </div>

    <div class="endpoint">
      <div class="endpoint-title">Funcionalidades</div>
      <ul>
        <li>Validação de CEP de 8 dígitos.</li>
        <li>Consulta de cidade utilizando a API do ViaCEP.</li>
        <li>Obtenção da temperatura em graus Celsius (°C), Fahrenheit (°F) e Kelvin (K).</li>
        <li>Tratamento de erros personalizados, retornando mensagens claras para o cliente.</li>
      </ul>
    </div>

    <div class="endpoint">
      <div class="endpoint-title">Resposta de Sucesso</div>
      <div class="example">
      <pre>
      {
        "city": "Mauá",
        "country": "Brasil",
        "localtime": "2025-01-27 10:00",
        "tempC": 25.0,
        "tempF": 77.0,
        "tempK": 298.0
      }
      </pre>
      </div>  
    </div>

    <h2>Resposta de Erro</h2>
    <div class="endpoint">
      <div class="endpoint-title">CEP Inválido</div>
      <p><strong>Status:</strong> <code>422 Unprocessable Entity</code></p>
      <div class="example">
        <pre>
        {
          "code": 422,
          "message": "invalid zipcode"
        }
        </pre>
      </div>

      <div class="endpoint-title">CEP Não Encontrado</div>
      <p><strong>Status:</strong> <code>404 Not Found</code></p>
      <div class="example">
        <pre>
        {
          "code": 404,
          "message": "can not find zipcode"
        }
        </pre>
      </div>  
    </div>
  </body>
</html>
`
	w.Write([]byte(html))
}
