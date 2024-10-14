package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/configs"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/event/handler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/infra/graph"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/infra/grpc/pb"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/infra/grpc/service"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/infra/web/webserver"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/pkg/events"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs.RabbitMQUrl)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		webOrderHandler.Create(w, r)
	})

	webserver.AddHandler("/list-orders",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			webOrderHandler.List(w, r)
		})

	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)
	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(rabbitMQUrl string) *amqp.Channel {
	conn, err := amqp.Dial(rabbitMQUrl)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
