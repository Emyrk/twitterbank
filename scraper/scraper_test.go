package scraper_test

import (
	"testing"

	"github.com/Emyrk/twitterbank/scraper"
)

func TestScraperInit(t *testing.T) {
	s, err := scraper.NewScraper("localhost", 8088)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	var _ = s
}
