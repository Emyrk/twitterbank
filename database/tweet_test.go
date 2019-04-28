package database_test

import (
	"testing"

	"github.com/Emyrk/twitterbank/database"
)

func TestDateString(t *testing.T) {
	dateString := "Thu Jan 24 15:01:51 +0000 2019"
	ts, err := database.ParseTwitterDate(dateString)
	if err != nil {
		t.Error(err)
	}
	if ts.Month() != 1 {
		t.Error("Wrong month")
	}
	if ts.Day() != 24 {
		t.Error("Wrong day")
	}
	if ts.Second() != 51 {
		t.Error("Wrong second")
	}
	if ts.Hour() != 15 {
		t.Error("Wrong hour")
	}
	var _ = ts
}
