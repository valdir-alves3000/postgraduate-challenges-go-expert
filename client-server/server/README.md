# Server

Este diretório contém o código para o servidor que fornece a cotação do dólar e a armazena em um banco de dados SQLite.

## Requisitos

- Go 1.18 ou superior
- SQLite (para o banco de dados)

## Instalação e Execução

1. **Instale as Dependências**:

   ```sh
   go mod tidy
   ```

2. **Compile e Execute o Servidor**:

   ```sh
   go run main.go
   ```

   O servidor estará disponível na porta `8080`.

## Como Funciona

- O servidor expõe um endpoint `/cotacao` que retorna a cotação do dólar em JSON.
- Ele obtém a cotação do dólar de uma API externa e a armazena em um banco de dados SQLite.

## Erros Comuns

- **"Database connection failed"**: Verifique se o SQLite está instalado e o caminho do banco de dados está correto.
- **"Failed to fetch quote"**: Verifique a conexão com a API externa.
- **"Failed to save quote to database"**: Verifique permissões e o estado do banco de dados.