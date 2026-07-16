package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"Goflix-Downloader/internal/models"
	"Goflix-Downloader/internal/queue"
)

func main() {
	fmt.Println("Goflix Producer: Feeding movies into RabbitMQ...")

	conn, ch, err := queue.ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("[-] Connection error: %v", err)
	}
	defer conn.Close()
	defer ch.Close()

	file, err := os.Open("data/movies.csv")
	if err != nil {
		log.Fatalf("[-] Failed to open CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("[-] Failed to read CSV: %v", err)
	}

	for _, record := range records {
		id, _ := strconv.Atoi(record[0])
		title := record[1]

		movie := models.Movie{
			ID:    id,
			Title: title,
		}

		err := queue.PublishMovie(ch, movie)
		if err != nil {
			fmt.Printf("[-] Failed to queue %s: %v\n", title, err)
		} else {
			fmt.Printf("[+] Movie queued -> #%d: %s\n", id, title)
		}
	}

	fmt.Println("All jobs published successfully!")
}