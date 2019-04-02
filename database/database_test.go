package database_test

import (
	"testing"

	"github.com/Emyrk/twitterbank/database"
)

func TestConfig(t *testing.T) {
	c := database.NewConfig(database.WithHost("Test"))
	if c.Host != "Test" {
		t.Errorf("Expected 'Test', got %s", c.Host)
	}
}
