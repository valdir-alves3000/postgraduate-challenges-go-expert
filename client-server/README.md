# Desafio de Cotação de Dólar com Go

Este projeto consiste em dois sistemas escritos em Go que trabalham juntos para obter e armazenar a cotação do dólar em relação ao real. O servidor expõe um endpoint HTTP para fornecer a cotação atual, e o cliente faz uma requisição a esse endpoint, salva a cotação em um arquivo e exibe uma mensagem de sucesso.

## Estrutura do Projeto

- ***[`client`](client)***: Cliente que faz a requisição ao servidor e salva a cotação em um arquivo.
- ***[`server`](server)***: Servidor que fornece a cotação do dólar e armazena os dados em um banco de dados SQLite.

## Requisitos

- Go 1.18 ou superior
- SQLite (para o servidor)
- Extensão SQLite no VS Code (opcional, para visualização do banco de dados)

## Instalação e Configuração

1. **Clone o Repositório**:

   ```sh
   git clone https://github.com/valdir-alves3000/postgraduate-challenges-go-expert/client-server.git
   cd client-server
   ```

2. **Instale as Dependências**:

   Navegue até o diretório do servidor e do cliente e execute:

   ```sh
   go mod tidy
   ```

3. **Compile e Execute o Servidor**:


   ```sh
   go run server/main.go
   ```

   O servidor estará disponível na porta `8080`.

4. **Compile e Execute o Cliente**:

   ```sh
   go run client/main.go
   ```

   O cliente fará uma requisição ao servidor, salvará a cotação em `cotacao.txt` e exibirá uma mensagem de sucesso.

## Como Funciona

### Servidor

O servidor expõe um endpoint `/cotacao` que retorna a cotação do dólar em relação ao real no formato JSON. Ele obtém a cotação de uma API externa e a armazena em um banco de dados SQLite.

### Cliente

O cliente faz uma requisição GET ao endpoint `/cotacao` do servidor. Recebe a cotação do dólar, salva essa cotação em um arquivo `cotacao.txt` e exibe uma mensagem indicando que a cotação foi salva.

## Exemplo de Requisição

Para testar o servidor, faça uma requisição GET ao endpoint `/cotacao`:

```sh
curl http://localhost:8080/cotacao
```

### Exemplo de Resposta

```json
{
  "code": "USD",
  "codein": "BRL",
  "name": "Dólar Americano/Real Brasileiro",
  "bid": "5.6238",
  "create_date": "2024-07-26 12:02:50"
}
```


## Erros Comuns

- **Para o Cliente**:
  - "Failed to get quote from server": Verifique se o servidor está em execução.
  - "Failed to save quote to file": Verifique permissões e o caminho do arquivo.

- **Para o Servidor**:
  - "Database connection failed": Verifique a instalação do SQLite.
  - "Failed to fetch quote": Verifique a conexão com a API externa.
  - "Failed to save quote to database": Verifique permissões e estado do banco de dados.
