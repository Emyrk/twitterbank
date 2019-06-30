package database

import (
	"fmt"
)

func (db *TwitterBankDatabase) InsertCompletedHeight(height int) error {
	if dbc := db.DB.Create(&CompletedHeight{BlockHeight: height}); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		return dbc.Error
	}
	return nil
}

func (db *TwitterBankDatabase) InsertNewUserChain(user *TwitterUser) error {
	if dbc := db.DB.Create(user); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		return dbc.Error
	}
	return nil
}

func (db *TwitterBankDatabase) InsertNewTweet(tweet *TwitterTweetObject, record *TwitterTweetRecord) error {
	// First check if the tweet exists. If it does, we only need to add another record of it's existence
	var exists Exists
	d := db.DB.Raw(`
		SELECT EXISTS(SELECT * FROM twitter_tweet_objects WHERE tweet_hash = ? AND tweet_id_str = ?)
	`, tweet.TweetHash, tweet.TweetIDStr).First(&exists)
	if d.Error != nil {
		return fmt.Errorf("exists_query: %s", d.Error.Error())
	}

	fmt.Printf("%v", tweet)

	tx := db.DB.Begin()
	if !exists.Exists {
		// Insert Tweet first
		if dbc := db.DB.Create(tweet); dbc.Error != nil {
			tx.Rollback()
			return fmt.Errorf("tweet_create: %s", dbc.Error.Error())
		}
	}
	if dbc := db.DB.Create(record); dbc.Error != nil {
		tx.Rollback()
		return fmt.Errorf("record_create: %s", dbc.Error)
	}

	if dbc := tx.Commit(); dbc.Error != nil {
		tx.Rollback()
		return fmt.Errorf("tx_commit: %s", dbc.Error)
	}
	return nil
}
