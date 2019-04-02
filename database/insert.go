package database

func (db *TwitterBankDatabase) InsertCompletedHeight(height int) error {
	if dbc := db.DB.Create(&Completed{BlockHeight: height}); dbc.Error != nil {
		// Create failed, do something e.g. return, panic etc.
		return dbc.Error
	}
	return nil
}
