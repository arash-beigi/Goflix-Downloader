package main

import (
	"Goflix-Downloader/internal/csv"
	"Goflix-Downloader/internal/downloader"
	"fmt"
)

func main() {
	fmt.Println("Starting Goflix Downloader System...")


	movieStream := csv.StreamMovies("data/movies.csv")


	downloader.StartDownloader(movieStream)

	fmt.Println(" All movies processed successfully!")
}