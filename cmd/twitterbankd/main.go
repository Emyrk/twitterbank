package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/Emyrk/twitterbank/scraper"

	"github.com/Emyrk/twitterbank/database"
	log "github.com/sirupsen/logrus"
)

type arrayFlags []string

var version = "v0.1.0"

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	enabledRoutines := arrayFlags{}

	var (
		factomdhost = flag.String("fhost", "localhost", "Factomd host")
		factomdport = flag.Int("fport", 8088, "Factomd port")

		postgreshost = flag.String("phost", "localhost", "Postgres host")
		postgresport = flag.Int("pport", 5432, "Postgres port")

		migrate = flag.Bool("m", false, "Automigrate on launch")
	)

	// For Debugging
	flag.Var(&enabledRoutines, "routine", "Can modify which routines are run")
	flag.Parse()

	db, err := database.NewDatabaseWithOpts(database.WithHost(*postgreshost), database.WithPort(*postgresport))
	if err != nil {
		panic(err)
	}
	log.Infof("Postgres database connected at %s:%d", *postgreshost, *postgresport)

	s, err := scraper.NewScaperWithDB(*factomdhost, *factomdport, db)
	if err != nil {
		panic(err)
	}

	// Graceful Close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			// sig is a ^C, handle it
			s.Close()
			return
		}
	}()

	log.Infof("Running Scraper %s", version)

	if len(enabledRoutines) == 0 {
		enabledRoutines = []string{"catchup"}
	}

	// Does as goroutine if not last
	do := func(f func(), i, l int) {
		if i == l {
			f()
		} else {
			go f()
		}
	}

	if *migrate {
		s.Database.AutoMigrate()
	}

	// Kinda hacky, but allows me to only run 1 routine if I want.
	for i, r := range enabledRoutines {
		switch r {
		case "catchup":
			do(s.Catchup, i, len(enabledRoutines)-1)
		}
	}
}
