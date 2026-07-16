package downloader

import (
	"Goflix-Downloader/internal/models"
	"Goflix-Downloader/internal/scraper"
	"fmt"
	"io"
	"net/http"
	"os"
)


func DownloadMovie(movie models.Movie) error {
	fmt.Printf("[WORKER] Searching for Movie #%d: %s on DonyayeSerial...\n", movie.ID, movie.Title)

	s := scraper.DonyaYeSerialScraper{}
	foundURL, err := s.Search(movie.Title)

	if err != nil || foundURL == "" {
		return fmt.Errorf("could not find any links for %s: %w", movie.Title, err)
	}

	fmt.Printf("[SUCCESS] Link found for %s -> %s\n", movie.Title, foundURL)

	savePath := fmt.Sprintf("data/downloads/%s.mp4", movie.Title)
	fmt.Printf("[DOWNLOAD] Starting download for %s...\n", movie.Title)

	err = downloadFileWithResume(foundURL, savePath)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}

	return nil
}

func downloadFileWithResume(url, savePath string) error {
	var startBytes int64 = 0
	fileInfo, err := os.Stat(savePath)
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
	var file *os.File
	if startBytes > 0 && resp.StatusCode == http.StatusPartialContent {
		file, err = os.OpenFile(savePath, os.O_WRONLY|os.O_APPEND, 0644)
	} else {
		file, err = os.Create(savePath)
	}
	if err != nil {
		return fmt.Errorf("failed to open/create destination file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write stream to file: %w", err)
	}

	return nil
}