package scraper

import (
	"fmt"
	"strings"
	"github.com/gocolly/colly/v2"
)

type DonyaYeSerialScraper struct{}

func (d DonyaYeSerialScraper) Search(title string) (string, error) {

	c := NewBaseCollector("donyayeserial.com")
	var downloadURL string

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasSuffix(link, ".mp4") || strings.HasSuffix(link, ".mkv") {
			downloadURL = link
		}
	})

	c.OnHTML(".post-title a[href], h2 a[href], h3 a[href]", func(e *colly.HTMLElement) {
		moviePageURL := e.Attr("href")
		
		if strings.Contains(moviePageURL, "donyayeserial.com/") && downloadURL == "" {
			e.Request.Visit(moviePageURL)
		}
	})

	searchURL := fmt.Sprintf("https://donyayeserial.com/?s=%s", strings.ReplaceAll(title, " ", "+"))

	err := c.Visit(searchURL)
	if err != nil {
		return "", err
	}

	if downloadURL == "" {
		return "", fmt.Errorf("not found on donyayeserial")
	}

	return downloadURL, nil
}