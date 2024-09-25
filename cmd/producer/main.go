package main

import (
	"context"
	"log"
	"time"

	"github.com/lokesh2201013/Microservices-go-rabbitmq/internal"
	ampq "github.com/rabbitmq/amqp091-go" // Define the alias ampq for the amqp091-go package
)

// handleError logs the error and exits the program if the error is non-nil.
func handleError(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err) // Log the error and exit the program.
	}
}

func main() {
	conn, err := internal.ConnectRabbitMQ("lokesh", "9910", "localhost:5672", "customers")
	handleError(err) // Use handleError to handle connection errors.
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	handleError(err) // Use handleError to handle client initialization errors.
	defer client.Close()

	err = client.CreateQueue("customers_created", true, false)
	handleError(err) // Use handleError to handle queue creation errors.

	err = client.CreateQueue("customers_test", false, true)
	handleError(err) // Use handleError to handle queue creation errors.

	err = client.CreateBinding("customers_created", "customer.created.*", "customer_events")
	handleError(err) // Use handleError to handle binding errors.

	err = client.CreateBinding("customers_test", "customer.*", "customer_events")
	handleError(err) // Use handleError to handle binding errors.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Send(ctx, "customer_events", "customers.created.us", ampq.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: ampq.Persistent,
		Body:         []byte(`A cool message between services`),
	})
	handleError(err) // Use handleError to handle message sending errors.

	err = client.Send(ctx, "customer_events", "customers.test", ampq.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: ampq.Transient,
		Body:         []byte(`An uncool message`),
	})
	handleError(err) // Use handleError to handle message sending errors.

	time.Sleep(10 * time.Second)

	log.Println(client)
}
