package csv

import (
	"Goflix-Downloader/internal/models"
	"encoding/csv"
	"io"
	"os"
	"strconv"
)

func StreamMovies(filepath string) <-chan models.Movie {
	ch := make(chan models.Movie)

	go func() {
		defer close(ch)

		file, err := os.Open(filepath)
		if err != nil {
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)

		for {
			record, err := reader.Read()

			if err == io.EOF {
				break
			}

			if err != nil {
				continue
			}

			id, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}

			movie := models.Movie{
				ID:    id,
				Title: record[1],
			}

			ch <- movie
		}
	}()

	return ch
}
