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
	var j = `{"":""}` //`{'Date_Recorded': '2019-06-26 09:11:15.869818', 'tweet': {'created_at': 'Wed Jun 26 13:11:10 +0000 2019', 'id': 1143869312240889856, 'id_str': '1143869312240889856', 'text': 'what up emyrk!', 'source': '<a href="http://twitter.com" rel="nofollow">Twitter Web Client</a>', 'truncated': False, 'in_reply_to_status_id': None, 'in_reply_to_status_id_str': None, 'in_reply_to_user_id': None, 'in_reply_to_user_id_str': None, 'in_reply_to_screen_name': None, 'user': {'id': 1128123860686200832, 'id_str': '1128123860686200832', 'name': 'FCT_Bot', 'screen_name': 'FCT_bot', 'location': None, 'url': None, 'description': None, 'translator_type': 'none', 'protected': False, 'verified': False, 'followers_count': 0, 'friends_count': 0, 'listed_count': 0, 'favourites_count': 0, 'statuses_count': 58, 'created_at': 'Tue May 14 02:24:22 +0000 2019', 'utc_offset': None, 'time_zone': None, 'geo_enabled': False, 'lang': None, 'contributors_enabled': False, 'is_translator': False, 'profile_background_color': 'F5F8FA', 'profile_background_image_url': '', 'profile_background_image_url_https': '', 'profile_background_tile': False, 'profile_link_color': '1DA1F2', 'profile_sidebar_border_color': 'C0DEED', 'profile_sidebar_fill_color': 'DDEEF6', 'profile_text_color': '333333', 'profile_use_background_image': True, 'profile_image_url': 'http://pbs.twimg.com/profile_images/1128124541845352449/UEshA5Ol_normal.png', 'profile_image_url_https': 'https://pbs.twimg.com/profile_images/1128124541845352449/UEshA5Ol_normal.png', 'default_profile': True, 'default_profile_image': False, 'following': None, 'follow_request_sent': None, 'notifications': None}, 'geo': None, 'coordinates': None, 'place': None, 'contributors': None, 'is_quote_status': False, 'quote_count': 0, 'reply_count': 0, 'retweet_count': 0, 'favorite_count': 0, 'entities': {'hashtags': [], 'urls': [], 'user_mentions': [], 'symbols': []}, 'favorited': False, 'retweeted': False, 'filter_level': 'low', 'lang': 'en', 'timestamp_ms': '1561554670759'}}`
	tweet, err := database.NewFactomTweetFromContent([]byte(j))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(tweet)
}
