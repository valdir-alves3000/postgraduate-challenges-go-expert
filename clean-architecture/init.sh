#!/bin/sh

until mysql -h "mysql" -u "root" -p"root" -e "SELECT 1" orders; do
    >&2 echo "Banco de dados não disponível, esperando..."
    sleep 5
done

echo "Banco de dados disponível. Executando migrações..."

migrate -path=sql/migrations -database "mysql://root:root@tcp(mysql:3306)/orders" -verbose up

echo "Iniciando o servidor..."

cd cmd/ordersystem
go run main.go wire_gen.go