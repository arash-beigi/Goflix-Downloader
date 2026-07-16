package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Goflix-Downloader/internal/downloader"
	"Goflix-Downloader/internal/models"
	"Goflix-Downloader/internal/queue"
)

func main() {
	fmt.Println(" Goflix Downloader: Waiting for movies (Max 5 concurrent downloads)...")

	conn, ch, err := queue.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("[-] Connection error: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	err = ch.Qos(5, 0, false)
	if err != nil {
		log.Fatalf("[-] QoS setup failed: %v", err)
	}

	msgs, err := ch.Consume(
		"movies_queue", 
		"",             
		false, 
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("[-] Failed to register consumer: %v", err)
	}

	sem := make(chan struct{}, 5)
	keepAlive := make(chan bool)

	go func() {
		for d := range msgs {
			sem <- struct{}{}

			go func(body []byte, deliveryTag uint64) {
				defer func() { <-sem }()

				var movie models.Movie
				err := json.Unmarshal(body, &movie)
				if err != nil {
					fmt.Printf("[-] Invalid message data: %v\n", err)
					ch.Nack(deliveryTag, false, false) 
					return
				}

				fmt.Printf("[WORKER] Started downloading: %s\n", movie.Title)

				err = downloader.DownloadMovie(movie)
				if err != nil {
					fmt.Printf("[-] Download failed for %s: %v\n", movie.Title, err)
					ch.Nack(deliveryTag, false, true) 
				} else {
					fmt.Printf("[SUCCESS] Download finished: %s\n", movie.Title)
					ch.Ack(deliveryTag, false) 
				}
			}(d.Body, d.DeliveryTag) 
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n Shutting down gracefully...")
		keepAlive <- true
	}()

	<-keepAlive
	fmt.Println("Downloader stopped.")
}