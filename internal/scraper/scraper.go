package scraper


type Scraper interface {
	Search(title string) (string, error)
}