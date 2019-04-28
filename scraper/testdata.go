package scraper

import (
	"fmt"
	"math/rand"

	"time"

	"github.com/Emyrk/twitterbank/database"
	"github.com/FactomProject/factomd/common/primitives"
)

func (s *Scraper) GenerateTestData() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		uid := rand.Int63()
		u := database.TwitterUser{
			UserID:    uid,
			UserIDStr: fmt.Sprintf("%d", uid),
			ChainID:   primitives.RandomHash().String(),
			EntryHash: primitives.RandomHash().String(),
		}
		s.Database.DB.Create(&u)

		for i := 0; i < 5; i++ {
			tid := rand.Int63()
			t := database.TwitterTweetObject{
				TweetID:          tid,
				TweetIDStr:       fmt.Sprintf("%d", tid),
				TweetAuthorID:    uid,
				TweetAuthorIDStr: fmt.Sprintf("%d", uid),
				TweetHash:        primitives.RandomHash().String(),
				EntryHash:        primitives.RandomHash().String(),
				ChainID:          primitives.RandomHash().String(),
				RawTweet:         ""}
			s.Database.DB.Create(&t)

			for i := 0; i < 5; i++ {
				r := database.TwitterTweetRecord{
					FactomRecorder:   primitives.RandomHash().String(),
					TweetID:          tid,
					TweetIDStr:       fmt.Sprintf("%d", tid),
					TweetAuthorID:    uid,
					TweetAuthorIDStr: fmt.Sprintf("%d", uid),
					TweetHash:        primitives.RandomHash().String(),
					Signature:        primitives.RandomHash().String() + primitives.RandomHash().String(),
					SigningKey:       primitives.RandomHash().String(),
					EntryHash:        primitives.RandomHash().String(),
					ChainID:          primitives.RandomHash().String(),
					BlockHeight:      rand.Intn(50000),
				}
				s.Database.DB.Create(&r)
			}
		}
	}
}
