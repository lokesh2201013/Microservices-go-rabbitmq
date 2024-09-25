package internal

import (
	"fmt"
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)


// RabbitClient encapsulates RabbitMQ connection and channel.
type RabbitClient struct {
	conn *amqp.Connection // Pointer to the RabbitMQ connection.
	ch   *amqp.Channel    // Pointer to the RabbitMQ channel.
}

// ConnectRabbitMQ establishes a connection to RabbitMQ with the given credentials.

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	// Format the URL for connecting to RabbitMQ using the provided credentials and vhost.
	url := fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost)
	// Connect to RabbitMQ and return the connection object.
	return amqp.Dial(url)
}

// NewRabbitMQClient initializes a new RabbitClient with a channel.
// `conn`: An established RabbitMQ connection.
// Returns a RabbitClient instance containing the connection and a channel.
func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	// Create a channel over the given connection.
	ch, err := conn.Channel()
	if err != nil {
		// Return an empty RabbitClient and an error if channel creation fails.
		return RabbitClient{}, err
	}
	

	// Return the initialized RabbitClient with the connection and channel.
	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

// Close closes the RabbitMQ channel and connection gracefully.
// Returns an error if closing either the channel or the connection fails.
func (rc RabbitClient) Close() error {
	// Close the channel.
	if err := rc.ch.Close(); err != nil {
		return err
	}
	// Close the connection.
	return rc.conn.Close()
}

// CreateQueue declares a new queue with the given name and properties.

func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(
		queueName,
		durable,  
		autoDelete,
		false,    
		false,  
		nil,       
	)

	return err
}

func (rc RabbitClient)CreateBinding(name,binding,exchange string)error{
	return  rc.ch.QueueBind(name,binding,exchange,false,nil)

}


func (rc RabbitClient) Send( ctx  context.Context, exchange ,routingKey string,options amqp.Publishing) error {
  return rc.ch.PublishWithContext(ctx,exchange,routingKey,true,false,options)
}
