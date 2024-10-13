//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/entity"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/event"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/infra/database"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/infra/web"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/internal/usecase"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/pkg/events"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(entity.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *usecase.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		usecase.NewCreateOrderUseCase,
	)
	return &usecase.CreateOrderUseCase{}
}

var setOrderUsecaseDependency = wire.NewSet(
	usecase.NewListOrdersUseCase,
)

func NewListOrdersUseCase(db *sql.DB) *usecase.ListOrdersUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		usecase.NewListOrdersUseCase,
	)
	return &usecase.ListOrdersUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		setOrderUsecaseDependency,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
