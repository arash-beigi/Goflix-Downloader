package downloader

import (
	"Goflix-Downloader/internal/models"
	"Goflix-Downloader/internal/scraper"
	"fmt"
	"sync"
)

func StartDownloader(movies <-chan models.Movie) {
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for movie := range movies {
		wg.Add(1)
		sem <- struct{}{}

		go func(m models.Movie) {
			defer wg.Done()
			defer func() { <-sem }()

			downloadMovie(m)
		}(movie)
	}

	wg.Wait()
	fmt.Println("All processes completed successfully!")
}

func downloadMovie(movie models.Movie) {
	fmt.Printf("[WORKER] Searching for Movie #%d: %s on DonyayeSerial...\n", movie.ID, movie.Title)


	s := scraper.DonyaYeSerialScraper{}
	foundURL, err := s.Search(movie.Title)

	if err != nil || foundURL == "" {
		fmt.Printf("[ERR] Could not find any links for %s. Error: %v\n", movie.Title, err)
		return
	}

	fmt.Printf("[SUCCESS] Link found for %s -> %s\n", movie.Title, foundURL)
}