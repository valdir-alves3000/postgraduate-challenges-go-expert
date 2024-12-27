# Rate Limiter com Go e Redis

Este projeto implementa um **Rate Limiter** configurável utilizando **Go** e **Redis**. Ele permite limitar requisições com base no **endereço IP** ou **token de acesso** e fornece um middleware que pode ser facilmente integrado em servidores HTTP.

## Funcionalidades

- **Limitação por IP:** Limita o número de requisições por segundo de um mesmo endereço IP.
- **Limitação por Token:** Permite configurar limites específicos para tokens de acesso únicos (passados no header `API_KEY`).
- **Bloqueio temporário:** Caso o limite seja excedido, o IP ou token é bloqueado por um período configurável.
- **Armazenamento no Redis:** Todas as informações do rate limiter são armazenadas no banco de dados Redis.
- **Configuração via `.env`:** O limite de requisições e o tempo de bloqueio são configuráveis através do arquivo `.env`.
- **Estratégia de Persistência:** O design permite a troca fácil do Redis por outros mecanismos de armazenamento.
- **Documentação da Rota `/docs`:** A aplicação inclui uma rota `/docs` com informações sobre o Rate Limiter.

---

## Pré-requisitos

- **Docker** e **Docker Compose** instalados.
- **Go 1.21+** instalado para desenvolvimento local.

---

## Instalação e Execução

1. **Clone o Repositório:**
   ```bash
   git clone https://github.com/valdir-alves3000/postgraduate-challenges-go-expert.git
   cd postgraduate-challenges-go-expert/rate-limiter
   ```

2. **Configure o Arquivo `.env`:**
   Exemplo de configuração:
   ```plaintext
   REDIS_HOST=redis
   REDIS_PORT=6379
   RATE_LIMIT=5               # Número máximo de requisições por segundo
   BLOCK_DURATION=300         # Duração do bloqueio em segundos
   SERVER_PORT=8080
   ```

3. **Suba o Docker Compose:**
   ```bash
   docker-compose up --build
   ```

4. **Testando a Aplicação:**
   - Acesse a rota **`/docs`** para visualizar a documentação.
   - Teste as rotas protegidas pelo Rate Limiter utilizando ferramentas como **Postman** ou **curl**.

---

## Rotas

- **`GET /docs`**: Documentação sobre a aplicação e como o Rate Limiter funciona.
- **Rotas Protegidas**: Qualquer rota implementada que utilize o **middleware** será protegida pelo Rate Limiter.

---

## **Exemplo de Teste**

- **Teste Limitação por IP**:
   ```bash
   curl http://localhost:8080/ -H "API_KEY: test"
   ```

   Faça múltiplas requisições consecutivas para verificar o bloqueio.

---

## Estratégia de Persistência

Por padrão, o sistema utiliza o **Redis** como mecanismo de persistência. Para trocar o Redis por outro sistema, basta implementar a interface definida em `internal/ratelimiter/interfaces.go` e registrar a nova estratégia.

---

## Testes

### **Testes de Integração**
Testam a eficácia do Rate Limiter com o Redis.

```bash
go test ./tests/integration
```

### Testes de Performance
Realizam testes sob carga para garantir robustez.

```bash
go test ./tests/performance
```

```bash
go test -bench=. ./tests/performance/ 
```

---

## Tecnologias Utilizadas

- **Go**: Linguagem de programação.
- **Redis**: Armazenamento de informações do Rate Limiter.
- **Docker** e **Docker Compose**: Facilita a configuração do ambiente.
- **Gorilla Mux**: Framework de roteamento HTTP.
- **Testing Framework**: Testes automatizados com o pacote padrão do Go.

---

## Como Funciona o Rate Limiter

1. **Requisições são interceptadas pelo middleware.**
2. O middleware verifica o **IP** ou **token** da requisição.
3. Caso o limite de requisições seja excedido:
   - A requisição é bloqueada.
   - Retorna **HTTP 429** com a mensagem:
     ```json
     {"error": "you have reached the maximum number of requests or actions allowed within a certain time frame"}
     ```
4. As informações são armazenadas no Redis com um tempo de expiração configurável.

