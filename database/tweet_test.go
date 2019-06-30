package database_test

import (
	"testing"

	"fmt"

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

func TestTweetParse(t *testing.T) {
	var j = `{"Date_Recorded": "2019-06-28 13:52:06.279026", "tweet": {"created_at": "Fri Jun 28 17:52:01 +0000 2019", "id": 1144664763840045058, "id_str": "1144664763840045058", "text": "test test test test test", "source": "<a href=\"http://twitter.com\" rel=\"nofollow\">Twitter Web Client</a>", "truncated": false, "in_reply_to_status_id": null, "in_reply_to_status_id_str": null, "in_reply_to_user_id": null, "in_reply_to_user_id_str": null, "in_reply_to_screen_name": null, "user": {"id": 1128123860686200832, "id_str": "1128123860686200832", "name": "FCT_Bot", "screen_name": "FCT_bot", "location": null, "url": null, "description": null, "translator_type": "none", "protected": false, "verified": false, "followers_count": 0, "friends_count": 0, "listed_count": 0, "favourites_count": 0, "statuses_count": 70, "created_at": "Tue May 14 02:24:22 +0000 2019", "utc_offset": null, "time_zone": null, "geo_enabled": false, "lang": null, "contributors_enabled": false, "is_translator": false, "profile_background_color": "F5F8FA", "profile_background_image_url": "", "profile_background_image_url_https": "", "profile_background_tile": false, "profile_link_color": "1DA1F2", "profile_sidebar_border_color": "C0DEED", "profile_sidebar_fill_color": "DDEEF6", "profile_text_color": "333333", "profile_use_background_image": true, "profile_image_url": "http://pbs.twimg.com/profile_images/1128124541845352449/UEshA5Ol_normal.png", "profile_image_url_https": "https://pbs.twimg.com/profile_images/1128124541845352449/UEshA5Ol_normal.png", "default_profile": true, "default_profile_image": false, "following": null, "follow_request_sent": null, "notifications": null}, "geo": null, "coordinates": null, "place": null, "contributors": null, "is_quote_status": false, "quote_count": 0, "reply_count": 0, "retweet_count": 0, "favorite_count": 0, "entities": {"hashtags": [], "urls": [], "user_mentions": [], "symbols": []}, "favorited": false, "retweeted": false, "filter_level": "low", "lang": "et", "timestamp_ms": "1561744321199"}}`

	tweet, err := database.NewFactomTweetFromContent([]byte(j))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tweet)
}
