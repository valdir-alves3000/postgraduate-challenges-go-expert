# Clean Architecture Go Application

Este projeto demonstra uma implementação de Clean Architecture em uma aplicação Go que interage com um banco de dados MySQL, RabbitMQ para mensagens e expõe uma API gRPC e GraphQL.

## Tabela de Conteúdos

- [Tecnologias](#tecnologias)
- [Configuração](#configuração)
- [Executando a Aplicação](#executando-a-aplicação)
- [Migrations](#migrations)
- [Endpoints da API](#endpoints-da-api)
- [Uso do Evans](#grpc)

## Tecnologias

- [**Go**](https://golang.org/): Linguagem de programação usada para a lógica da aplicação.
- [**MySQL**](https://www.mysql.com/): Banco de dados para armazenar pedidos.
- [**RabbitMQ**](https://www.rabbitmq.com/): Broker de mensagens para lidar com eventos.
- [**gRPC**](https://grpc.io/): Framework de chamada de procedimento remoto para comunicação de serviços.
- [**GraphQL**](https://graphql.org/): Linguagem de consulta para APIs.
- [**Docker**](https://www.docker.com/): Containerização para configuração do ambiente.
- [**Golang-Migrate**](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate): Ferramenta de migração de banco de dados.
- [**Evans**](https://github.com/ktr0731/evans): Cliente CLI para interagir com APIs gRPC.
- [**Wire**](https://github.com/google/wire): Ferramenta de injeção de dependências para Go.

## Configuração

Certifique-se de ter o Docker instalado em sua máquina.

### Clonar o Repositório

```bash
git clone https://github.com/valdir-alves3000/postgraduate-challenges-go-expert.git
cd postgraduate-challenges-go-expert/clean-architecture
```

## Executando a Aplicação

1. Inicie os serviços usando o Docker Compose:

   ```bash
   docker-compose up -d
   ```

2. A aplicação Go será iniciada automaticamente ao executar os serviços do Docker.

## Migrations

### Instalando a Ferramenta de Migrations

Para gerenciar as migrações do banco de dados MySQL, você precisará instalar a ferramenta `golang-migrate`. Execute o seguinte comando para instalá-la:

```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Executando as Migrations

Após instalar a ferramenta de migrations, você pode executar as migrations no banco de dados MySQL utilizando o seguinte comando:

```bash
migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up
```

Certifique-se de que o banco de dados MySQL esteja rodando e que o diretório `sql/migrations` contenha os arquivos de migração.

### Iniciando a Aplicação

Após clonar o repositório, siga estas etapas para garantir que todas as dependências do projeto estejam atualizadas e a aplicação seja executada corretamente:

1. **Atualizar Dependências**  
   Execute o seguinte comando para organizar e garantir que todas as dependências do projeto estejam atualizadas:

   ```bash
   go mod tidy
   ```

2. **Navegar até o Diretório do Projeto**  
   Mude para o diretório onde o código da aplicação está localizado:

   ```bash
   cd cmd/ordersystem
   ```

3. **Executar a Aplicação**  
   Inicie a aplicação com o comando a seguir:

   ```bash
   go run main.go wire_gen.go
   ```

## Endpoints da API

### Endpoints Web

- **Criar Pedido**
  - **Método**: `POST`
  - **URL**: `http://localhost:8000/order`
  - **Corpo da Requisição**:
    ```json
    {
        "id": "order-ID",
        "price": 100.5,
        "tax": 0.5
    }
    ```

- **Listar Pedidos**
  - **Método**: `GET`
  - **URL**: `http://localhost:8000/list-orders`

### GraphQL

- **Mutação para Criar Pedido**
  ```graphql
  mutation createOrder {
    createOrder(input: {id: "order-ID", Price: 162, Tax: 2}) {
      id
      Price
      Tax
      FinalPrice
    }
  }
  ```

- **Consulta para Listar Pedidos**
  ```graphql
  query ListOrders {
    listOrders {
      id
      Price
      Tax
      FinalPrice
    }
  }
  ```

### gRPC

#### Instalando o Evans

Para instalar o Evans, que será usado para interagir com a API gRPC, execute o seguinte comando:

```bash
go install github.com/ktr0731/evans@latest
```

- **Evans**  
  Para interagir com o serviço gRPC utilizando o Evans, execute os seguintes comandos:

  ```bash
  evans -r repl
  ```

  Após iniciar o REPL do Evans, você pode acessar o serviço `OrderService` da seguinte forma:

  ```plaintext
  127.0.0.1:50051> package pb
  pb@127.0.0.1:50051> service OrderService
  ```

  Com isso, você estará preparado para fazer chamadas aos métodos `CreateOrder` e `ListOrders` disponíveis no `OrderService`.

#### Exemplo de `CreateOrder`

Para criar um pedido, você pode chamar o método `CreateOrder` assim:

```plaintext
pb.OrderService@127.0.0.1:50051> call CreateOrder
```

Quando solicitado, forneça os parâmetros do pedido no formato JSON:

```plaintext
{
  "id": "order-1",
  "price": 100.5,
  "tax": 0.5
}
```

#### Exemplo de `ListOrders`

Para listar todos os pedidos, você pode chamar o método `ListOrders` assim:

```plaintext
pb.OrderService@127.0.0.1:50051> call ListOrders
```

O `ListOrders` não requer parâmetros e retornará uma lista de todos os pedidos registrados.