package scraper

import (
	"fmt"

	"time"

	"github.com/Emyrk/factom-raw"
	"github.com/Emyrk/twitterbank/database"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

var scraperlog = log.WithFields(log.Fields{"file": "scraper.go"})

type Scraper struct {
	Factom   factom_raw.Fetcher
	Database *database.TwitterBankDatabase
	quit     bool
}

type DatabaseConfig struct {
	Host     string
	Port     int
	DBName   string
	Password string
}

func NewScraper(host string, port int, dbopts ...func(config *database.TwitterBankConfig)) (*Scraper, error) {
	flog := scraperlog.WithField("func", "NewScraper")

	s := new(Scraper)
	factomd := fmt.Sprintf("%s:%d", host, port)
	s.Factom = factom_raw.NewAPIReader(factomd)
	_, err := s.Factom.FetchDBlockHead()
	if err != nil {
		return nil, err
	}
	flog.Infof("Factomd location %s", factomd)

	db, err := database.NewDatabaseWithOpts(dbopts...)
	if err != nil {
		panic(err)
	}
	s.Database = db

	flog.Infof("Postgres database connected")

	return s, nil
}

func (s *Scraper) Close() {
	s.quit = true
}

var CurrentCatchup uint32
var CurrentTop uint32

// Catchup will sync all entries in order. It will sync full block heights, not partials
// so the process can be stopped and restarted at any point.
func (s *Scraper) Catchup() {
	flog := scraperlog.WithFields(log.Fields{"func": "CatchUp"})
	flog.Info("Catchup started")
	// Find the highest height completed
	comp, err := s.Database.FetchHighestDBInserted()
	if err != nil {
		panic(err)
	}
	next := uint32(comp + 1)

	getNextTop := func() uint32 {
		for {
			topDblock, err := s.Factom.FetchDBlockHead()
			if err != nil {
				flog.Error(err)
				time.Sleep(3 * time.Second)
				continue
			}
			return topDblock.GetDatabaseHeight()
		}
	}

	start := time.Now()
	top := getNextTop()
	CurrentTop = top

MainCatchupLoop:
	for {
		if s.quit {
			s.Database.Close()
			return
		}

		if next%10 == 0 {
			flog.WithFields(log.Fields{"current": next, "top": top, "time": time.Since(start)}).Info("")
		}
		start = time.Now()
		if next > top {
			top = getNextTop()
			if next > top {
				time.Sleep(30 * time.Second)
				continue
			}
		}
		CurrentCatchup = next

		dblock, err := s.Factom.FetchDBlockByHeight(next)
		if err != nil {
			errorAndWait(flog.WithField("fetch", "dblock"), err)
			continue
		}

		fblock, err := s.Factom.FetchFBlockByHeight(next)
		if err != nil {
			errorAndWait(flog.WithField("fetch", "fblock"), err)
			continue
		}
		ablock, err := s.Factom.FetchABlockByHeight(next)
		if err != nil {
			errorAndWait(flog.WithField("fetch", "ablock"), err)
			continue
		}

		ec_keymr := dblock.GetDBEntries()[1]
		ecblock, err := s.Factom.FetchECBlock(ec_keymr.GetKeyMR())
		if err != nil {
			errorAndWait(flog.WithField("fetch", "ecblock"), err)
			continue
		}

		var _, _, _ = fblock, ablock, ecblock

		err = s.Database.InsertCompletedHeight(int(next))
		if err != nil {
			errorAndWait(flog.WithField("insert", "completed"), err)
			continue MainCatchupLoop
		}
		// End loop
		next++
	}
}

func errorAndWait(logger *log.Entry, err error) {
	logger.Error(err)
	time.Sleep(2 * time.Second)
}
