FROM golang:latest

WORKDIR /app
RUN apt-get update && apt-get install -y default-mysql-client
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
COPY cmd/ordersystem/.env ./cmd/ordersystem/

RUN go mod tidy

RUN chmod +x ./cmd/ordersystem/.env

RUN chmod +x init.sh

CMD ["./init.sh"]
