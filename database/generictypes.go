package database

type Exists struct {
	Exists bool `gorm: exists`
}
