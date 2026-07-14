package scraper

import (
	"time"

	"github.com/gocolly/colly/v2"
)

func NewBaseCollector(domain string) *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains(domain, "www."+domain),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*" + domain + "*",
		Parallelism: 1,
		Delay:       2 * time.Second,
		RandomDelay: 2 * time.Second,
	})
	return c
}
