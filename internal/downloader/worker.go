package downloader

import (
	"Goflix-Downloader/internal/models"
	"Goflix-Downloader/internal/scraper"
	"fmt"
	"net/http"
	"os"
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


func downloadFileWithResume (url ,savaPath string) error {
	var startBytes int64 = 0
	fileInfo , err := os.Stat(savaPath)
	if err == nil {
		startBytes = fileInfo.Size()
		fmt.Printf("[INFO] Partial file found (%d bytes). Resuming download...\n", startBytes)
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
			return fmt.Errorf("failed to create http request: %w", err)
	}

	if startBytes > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startBytes))
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("network request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("server responded with bad status: %s", resp.Status)
	}
}