package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

// CompletedHeight indicates what heights have already been synced
type CompletedHeight struct {
	gorm.Model
	BlockHeight int
}

/********
	Twitter Objects
 ********/

// TwitterTweetObject is the object per tweet
//	Does not include all objects for tweets.
type TwitterTweetObject struct {
	//gorm.Model
	// Integer based indexes are faster. Twitter reccomends the string
	TweetAuthorID    int64  `json:"tweet_author"`
	TweetAuthorIDStr string `json:"tweet_author_str"`
	ChainID          string `json:"chain_id"`
	EntryHash        string `json:"entry_hash"`

	TweetID    int64  `json:"tweet_id"`
	TweetIDStr string `json:"tweet_id_str" gorm:"primary_key"`
	TweetHash  string `json:"tweet_hash" gorm:"primary_key"`

	// The full raw tweet reponse
	RawTweet string `json:"raw_tweet"`

	// Some fields will be parsed for better searching
	//	The time the tweet was tweeted
	TweetCreatedAt time.Time `json:"tweeted_time"`

	// TODO: Handle Quotes and Retweets
}

// Collaborate checks the fields in the tweet object match that in the content.
//	One is built from the header, one from the content.
func (t TwitterTweetObject) Collaborate(content *FactomTweetContent) bool {
	tu := content.Tweet.User
	if tu.IDStr != t.TweetAuthorIDStr {
		return false
	}

	if content.Tweet.IDStr != t.TweetIDStr {
		return false
	}

	// TODO: Probably check more
	return true
}

type TwitterTweetRecord struct {
	// Factom Identity
	FactomRecorder string `json:"record_identity" gorm:"primary_key"`

	// Twitter Related IDs
	TweetAuthorID    int64  `json:"tweet_author"`
	TweetAuthorIDStr string `json:"tweet_author_str"`
	TweetID          int64  `json:"tweet_id"`
	TweetIDStr       string `json:"tweet_id_str" gorm:"primary_key"`
	TweetHash        string `json:"tweet_hash"`

	ChainID    string `json:"chain_id"`
	EntryHash  string `json:"entry_hash"`
	SigningKey string `json:"signing_key"`
	Signature  string `json:"signature"`

	//	The time recorded into factom
	TweetRecordedAt time.Time `json:"tweet_recorded_time"`
	BlockHeight     int       `json:"block_height"`
}

type TwitterUser struct {
	// Integer based indexes are faster. Twitter reccomends the string
	UserID            int64  `json:"user_id"`
	UserIDStr         string `json:"user_id_str" gorm:"primary_key"`
	BlockInitiated    int    `json:"block_initiated"`      // Block chain was started
	BlockInitatedUnix int64  `json:"block_initiated_unix"` // Block chain was started
	ChainID           string `json:"chain_id" gorm:"unique_key"`
	EntryHash         string `json:"entry_hash"`

	// Associations
	Tweets []TwitterTweetObject `json:"tweets,omitempty" gorm:"foreignkey:TweetAuthorIDStr;association_foreignkey:UserIDStr"`
}

func (u *TwitterUser) FindTweets(db *gorm.DB, limit, offset int) error {
	dbc := db.Limit(limit).Offset(offset).Model(u).Related(&u.Tweets, "TweetAuthorIDStr")
	return dbc.Error
}

func (db *TwitterBankDatabase) AutoMigrate() {
	panicDBErr(db.DB.AutoMigrate(&CompletedHeight{}))
	panicDBErr(db.DB.AutoMigrate(&TwitterTweetObject{}))
	panicDBErr(db.DB.AutoMigrate(&TwitterTweetRecord{}))
	panicDBErr(db.DB.AutoMigrate(&TwitterUser{}))
}

func panicDBErr(db *gorm.DB) {
	if db.Error != nil {
		panic(db.Error)
	}
}
