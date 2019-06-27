package database

import "github.com/jinzhu/gorm"

func (db *TwitterBankDatabase) FetchHighestDBInserted() (int, error) {
	c := new(CompletedHeight)
	if dbc := db.DB.Last(c); dbc.Error != nil {
		if dbc.Error == gorm.ErrRecordNotFound {
			return -1, nil
		}
		// Create failed, do something e.g. return, panic etc.
		return -2, dbc.Error
	}
	return c.BlockHeight, nil
}

func (db *TwitterBankDatabase) FetchUserByUID(uid string) (*TwitterUser, error) {
	c := TwitterUser{}
	dbc := db.DB.Where("user_id_str = ?", uid).Find(&c)
	if dbc.Error != nil {
		return nil, dbc.Error
	}
	return &c, nil
}

func (db *TwitterBankDatabase) FetchAllUsers() ([]TwitterUser, error) {
	users := []TwitterUser{}
	dbc := db.DB.Find(&users)
	if dbc.Error != nil {
		return nil, dbc.Error
	}
	return users, nil
}

func (db *TwitterBankDatabase) FetchTotalNumberOfUsers() (int, error) {
	return db.TableCount("twitter_users")

}

func (db *TwitterBankDatabase) FetchTotalNumberOfTweets() (int, error) {
	return db.TableCount("twitter_tweet_objects")

}

func (db *TwitterBankDatabase) FetchTotalNumberOfTweetRecords() (int, error) {
	return db.TableCount("twitter_tweet_records")
}

func (db *TwitterBankDatabase) TableCount(table string) (int, error) {
	count := 0
	dbc := db.DB.Table(table).Count(&count)
	if dbc.Error != nil {
		return count, dbc.Error
	}
	return count, nil
}

func (db *TwitterBankDatabase) FetchTweetByTID(tid string) (*TwitterTweetObject, error) {
	c := TwitterTweetObject{}
	dbc := db.DB.Where("tweet_id_str = ?", tid).Find(&c)
	if dbc.Error != nil {
		return nil, dbc.Error
	}
	return &c, nil
}
