package database

import "github.com/jinzhu/gorm"

// Completed indicates what heights have already been synced
type Completed struct {
	gorm.Model
	BlockHeight int
}

func (db *TwitterBankDatabase) AutoMigrate() {
	db.DB.AutoMigrate(&Completed{})
}
