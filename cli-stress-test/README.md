# CLI Stress Test

Uma ferramenta de linha de comando para teste de carga em serviços web, desenvolvida em Go.

## Funcionalidades
- Execução concorrente de requisições.
- Configuração personalizável do número total de requisições.
- Ajuste do nível de concorrência.
- Métricas detalhadas de desempenho.
- Análise da distribuição de códigos de status HTTP.

---

## Instalação

### 1. **Clonar o repositório**
Execute os seguintes comandos no terminal:

```bash
git clone https://github.com/valdir-alves3000/postgraduate-challenges-go-expert.git

cd postgraduate-challenges-go-expert/cli-stress-test
```

### 2. **Build local**
Gere o executável da aplicação com:

```bash
go build -o cli-stress-test ./cmd/main.go
```

Se o Go estiver configurado corretamente, você terá o binário `cli-stress-test` pronto para uso.

---

## Execução

### Opções de linha de comando
- `--url`: URL de destino para o teste (obrigatório).
- `--requests`: Número total de requisições a serem feitas.
- `--concurrency`: Número de requisições concorrentes.

### Exemplos

#### Uso básico

```bash
./cli-stress-test --url=https://api.github.com/users --requests=1000 --concurrency=10
```

#### Docker
Se preferir executar via Docker:

1. **Build da imagem Docker:**

```bash
docker build -t cli-stress-test .
```

2. **Executar a imagem:**

```bash
docker run --rm cli-stress-test --url=https://api.github.com/users --requests=1000 --concurrency=10
```

#### Docker Hub
Para usar a imagem disponível no Docker Hub:

```bash
docker run --rm valdiralves3000/cli-stress-test --url=https://api.github.com/users --requests=1000 --concurrency=10
```

---

## Saída
A ferramenta gera um relatório contendo:
- Tempo total de execução.
- Número total de requisições realizadas.
- Contagem de requisições bem-sucedidas (HTTP 200).
- Distribuição dos códigos de status HTTP.