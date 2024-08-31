# CEP Lookup - Go API Client

Este projeto é um cliente em Go para busca de endereços através de CEPs (Códigos de Endereçamento Postal) utilizando duas APIs brasileiras: BrasilAPI e ViaCEP. O programa faz requisições paralelas para ambas as APIs e retorna a resposta da que responder mais rápido.

## Funcionalidades

- Consulta de endereços a partir de CEPs.
- Uso de múltiplas APIs para maior confiabilidade.
- Tratamento de timeout para garantir que o programa não fique preso em uma API lenta.

## Requisitos

- Go 1.19 ou superior.

## Instalação

1. Clone o repositório:

   ```sh
   git clone https://github.com/valdir-alves3000/postgraduate-challenges-go-expert/multithreading.git
   cd multithreading/cep-lookup
   ```

## Uso

Você pode executar o programa diretamente usando o comando `go run` e passando um ou mais CEPs como argumentos:

```sh
go run main.go 01001000 02002000 03003000
```

O programa retornará o endereço correspondente ao CEP mais rápido encontrado pelas APIs.

### Exemplo de Saída

```sh
  ========================================
Resultado mais rápido vindo de brasilapi.com.br:
CEP: 01001-000
Logradouro: Praça da Sé
Bairro: Sé
Localidade: São Paulo
UF: SP
```

Caso ocorra um timeout:

```sh
Erro: Timeout - Nenhuma das APIs respondeu em 1.0 segundo.
```

## Estrutura do Código

- **handlerAddress**: Função responsável por realizar a requisição HTTP para uma das APIs e retornar o endereço encontrado.
- **main**: Função principal que gerencia as requisições para os CEPs fornecidos, tratando a resposta mais rápida.
