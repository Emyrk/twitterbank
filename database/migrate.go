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
	TweetID    int64  `json:"tweet_id" gorm:"unique_index"`
	TweetIDStr string `json:"tweet_id_str" gorm:"primary_key"`

	// The full raw tweet reponse
	RawTweet string `json:"raw_tweet"`

	// Some fields will be parsed for better searching
	//	The time the tweet was tweeted
	TweetCreatedAt time.Time `json:"tweet_created_time"`
	//	The time recorded into factom
	TweetRecordedAt  time.Time `json:"tweet_recorded_time"`
	TweetAuthorID    int64     `json:"tweet_author"`
	TweetAuthorIDStr string    `json:"tweet_author_str"`

	// TODO: Handle Quotes and Retweets
}

type TwitterUser struct {
	// Integer based indexes are faster. Twitter reccomends the string
	UserID    int64  `json:"user_id" gorm:"unique_index"`
	UserIDStr string `json:"user_id_str" gorm:"primary_key"`

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
	panicDBErr(db.DB.AutoMigrate(&TwitterUser{}))
}

func panicDBErr(db *gorm.DB) {
	if db.Error != nil {
		panic(db.Error)
	}
}
