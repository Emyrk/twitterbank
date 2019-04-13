package scraper

import (
	"fmt"
	"math/rand"

	"github.com/Emyrk/twitterbank/database"
)

func (s *Scraper) GenerateTestData() {
	for i := 0; i < 20; i++ {
		uid := rand.Int63()
		u := database.TwitterUser{UserID: uid, UserIDStr: fmt.Sprintf("%d", uid)}
		s.Database.DB.Create(&u)

		for i := 0; i < 25; i++ {
			tid := rand.Int63()
			t := database.TwitterTweetObject{TweetID: tid, TweetIDStr: fmt.Sprintf("%d", tid), TweetAuthorID: uid, TweetAuthorIDStr: fmt.Sprintf("%d", uid)}
			s.Database.DB.Create(&t)
		}
	}
}
