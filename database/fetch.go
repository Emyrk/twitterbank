package database

import "github.com/jinzhu/gorm"

func (db *TwitterBankDatabase) FetchHighestDBInserted() (int, error) {
	c := new(Completed)
	if dbc := db.DB.Last(c); dbc.Error != nil {
		if dbc.Error == gorm.ErrRecordNotFound {
			return -1, nil
		}
		// Create failed, do something e.g. return, panic etc.
		return -2, dbc.Error
	}
	return c.BlockHeight, nil
}
