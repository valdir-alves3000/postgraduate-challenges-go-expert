# Auction Go

Este projeto simula a implementação de um sistema de leilões com uma funcionalidade de fechamento automático de leilões expirados. O código foi desenvolvido em Go e inclui uma camada de testes unitários para validar a lógica do sistema. O projeto também pode ser facilmente executado utilizando Docker.

## Inicialização com Docker

Para rodar o projeto utilizando Docker, siga os passos abaixo:

### 1. Clone o repositório

Primeiro, clone o repositório do GitHub:

```bash
git clone https://github.com/valdir-alves3000/postgraduate-challenges-go-expert.git
cd labs-auction-goexpert
```

### 2. Construindo e rodando com Docker

Certifique-se de ter o Docker instalado na sua máquina.

No diretório raiz do projeto, você encontrará um arquivo `docker-compose.yml` que define dois serviços: a aplicação e o MongoDB.

Para construir e iniciar os serviços, execute o seguinte comando:

```bash
docker-compose up --build
```

Isso irá construir a imagem do Go e iniciar tanto a aplicação quanto o MongoDB.

### 3. Acessando os serviços

A aplicação estará disponível na porta `8080` e o MongoDB na porta `27017`.

- Acesse a API na URL: [http://localhost:8080](http://localhost:8080)
- O MongoDB estará disponível na URL: [mongodb://localhost:27017](mongodb://localhost:27017)

### 4. Parando os serviços

Para parar os serviços, execute:

```bash
docker-compose down
```

## Exemplos de Rotas

Este projeto simula o gerenciamento de leilões. Algumas rotas de exemplo são:

### 1. **Criar Leilão**

`POST /auctions`

```bash
", 
    "condition": 1, 
    "status": 0
    }' \curl -X POST http://localhost:8080/auctions \
  -d '{
    "product_name": "Product A", 
    "category": "Category A", 
    "description": "Auction Description
  -H "Content-Type: application/json"
```

### 2. **Buscar Leilão por ID**

`GET /auctions/{id}`

```bash
curl http://localhost:8080/auctions/1
```

### Mais rotas no diretorio [auction.http](api/auction.http)

### Rodando Testes no Arquivo `start_auctions_closer_test.go`

Caso deseje rodar o teste do arquivo `start_auctions_closer_test.go`, execute:

```bash
go test ./cmd/auction/start_auctions_closer_test.go
```

Este arquivo de teste contém a lógica de fechamento automático de leilões expirados.