# Client

Este diretório contém o código para o cliente que se conecta ao servidor e salva a cotação do dólar em um arquivo.


## Requisitos

- Go 1.18 ou superior

## Instalação e Execução

1. **Instale as Dependências**:

   Navegue até o diretório `client` e execute:

   ```sh
   go mod tidy
   ```

2. **Compile e Execute o Cliente**:

   No diretório `client`, execute:

   ```sh
   go run main.go
   ```

   O cliente fará uma requisição ao servidor para obter a cotação do dólar e salvará a cotação em um arquivo chamado `cotacao.txt`.

## Como Funciona

- O cliente faz uma requisição GET ao endpoint `/cotacao` do servidor.
- Recebe a cotação do dólar e salva essa informação em `cotacao.txt`.

## Erros Comuns

- **"Failed to get quote from server"**: Verifique se o servidor está em execução e acessível.
- **"Failed to save quote to file"**: Verifique permissões e o caminho do arquivo.