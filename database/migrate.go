package database

import "github.com/jinzhu/gorm"

// CompletedHeight indicates what heights have already been synced
type CompletedHeight struct {
	gorm.Model
	BlockHeight int
}

func (db *TwitterBankDatabase) AutoMigrate() {
	db.DB.AutoMigrate(&CompletedHeight{})
}
