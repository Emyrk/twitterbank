package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type TwitterBankDatabase struct {
	DB *gorm.DB
}

func NewDatabaseWithOpts(opts ...func(*TwitterBankConfig)) (*TwitterBankDatabase, error) {
	c := NewConfig(opts...)
	return NewDatabase(c)
}

func NewDatabase(c *TwitterBankConfig) (*TwitterBankDatabase, error) {
	tb := new(TwitterBankDatabase)
	if c == nil {
		c = NewConfig()
	}

	db, err := gorm.Open(c.DBDialect, c.ConnectString())
	if err != nil {
		return nil, err
	}

	tb.DB = db

	return tb, nil
}

func (db *TwitterBankDatabase) Close() error {
	return db.DB.Close()
}

type TwitterBankConfig struct {
	Host      string
	Password  string
	Port      int
	User      string
	DBName    string
	DBDialect string // 'postgres'
	SSL       bool
}

func (t TwitterBankConfig) ConnectString() string {
	additional := ""
	if !t.SSL {
		additional += " sslmode=disable"
	}
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s %s", t.Host, t.Port, t.User, t.DBName, t.Password, additional)
}

func NewConfig(opts ...func(*TwitterBankConfig)) *TwitterBankConfig {
	c := new(TwitterBankConfig)
	c.Host = "localhost"
	c.Password = "password"
	c.Port = 5432
	c.User = "postgres"
	c.DBName = "postgres"
	c.DBDialect = "postgres"
	c.SSL = false

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithHost(host string) func(*TwitterBankConfig) {
	return func(c *TwitterBankConfig) {
		c.Host = host
	}
}

func WithPort(port int) func(*TwitterBankConfig) {
	return func(c *TwitterBankConfig) {
		c.Port = port
	}
}

func WithPassword(password string) func(*TwitterBankConfig) {
	return func(c *TwitterBankConfig) {
		c.Password = password
	}
}
