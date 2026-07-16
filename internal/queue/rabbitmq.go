package queue



import (
	"encoding/json"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"Goflix-Downloader/internal/models"


)

func ConnectRabbitMQ()(*amqp091.Connection, *amqp091.Channel, error){
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to open channel: %w", err)
	}
	_, err = ch.QueueDeclare(
		"movies_queue",
		true,           
		false,          
		false,          
		false,          
		nil,            
	)

	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return conn, ch, nil
}

func PublishMovie(ch *amqp091.Channel, movie models.Movie) error {
	body, err := json.Marshal(movie)
	if err != nil {
		return fmt.Errorf("failed to marshal movie: %w", err)
	}

	return ch.Publish(
		"",             
		"movies_queue", 
		false,          
		false,          
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp091.Persistent, 
		},
	)
}
	
