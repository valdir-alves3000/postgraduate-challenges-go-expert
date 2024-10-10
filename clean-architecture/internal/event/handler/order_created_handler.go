package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
	"github.com/valdir-alves3000/postgraduate-challenges-go-expert/clean-architecture/pkg/events"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Order created: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	queue, err := h.RabbitMQChannel.QueueDeclare(
		"orders", // queue name
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		fmt.Println("Failed to declare the queue:", err)
		return
	}

	h.RabbitMQChannel.Publish(
		"",          // exchange
		queue.Name,  // queue name (bind to "orders")
		false,       // mandatory
		false,       // immediate
		msgRabbitmq, // message to publish
	)
}
