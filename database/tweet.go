package database

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

type FactomTweetContent struct {
	DateRecorded string        `json:"Date_Recorded"`
	Tweet        twitter.Tweet `json:tweet`
	// Tweet        struct {
	// 	CreatedAt            string `json:"created_at"`
	// 	ID                   int64  `json:"id"`
	// 	IDStr                string `json:"id_str"`
	// 	Text                 string `json:"text"`
	// 	Source               string `json:"source"`
	// 	Truncated            bool   `json:"truncated"`
	// 	InReplyToStatusID    string `json:"in_reply_to_status_id"`
	// 	InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
	// 	InReplyToUserID      string `json:"in_reply_to_user_id"`
	// 	InReplyToUserIDStr   string `json:"in_reply_to_user_id_str"`
	// 	InReplyToScreenName  string `json:"in_reply_to_screen_name"`
	// 	User                 struct {
	// 		ID                             int64  `json:"id"`
	// 		IDStr                          string `json:"id_str"`
	// 		Name                           string `json:"name"`
	// 		ScreenName                     string `json:"screen_name"`
	// 		Location                       string `json:"location"`
	// 		URL                            string `json:"url"`
	// 		Description                    string `json:"description"`
	// 		TranslatorType                 string `json:"translator_type"`
	// 		Protected                      bool   `json:"protected"`
	// 		Verified                       bool   `json:"verified"`
	// 		FollowersCount                 int    `json:"followers_count"`
	// 		FriendsCount                   int    `json:"friends_count"`
	// 		ListedCount                    int    `json:"listed_count"`
	// 		FavouritesCount                int    `json:"favourites_count"`
	// 		StatusesCount                  int    `json:"statuses_count"`
	// 		CreatedAt                      string `json:"created_at"`
	// 		UtcOffset                      string `json:"utc_offset"`
	// 		TimeZone                       string `json:"time_zone"`
	// 		GeoEnabled                     bool   `json:"geo_enabled"`
	// 		Lang                           string `json:"lang"`
	// 		ContributorsEnabled            bool   `json:"contributors_enabled"`
	// 		IsTranslator                   bool   `json:"is_translator"`
	// 		ProfileBackgroundColor         string `json:"profile_background_color"`
	// 		ProfileBackgroundImageURL      string `json:"profile_background_image_url"`
	// 		ProfileBackgroundImageURLHTTPS string `json:"profile_background_image_url_https"`
	// 		ProfileBackgroundTile          bool   `json:"profile_background_tile"`
	// 		ProfileLinkColor               string `json:"profile_link_color"`
	// 		ProfileSidebarBorderColor      string `json:"profile_sidebar_border_color"`
	// 		ProfileSidebarFillColor        string `json:"profile_sidebar_fill_color"`
	// 		ProfileTextColor               string `json:"profile_text_color"`
	// 		ProfileUseBackgroundImage      bool   `json:"profile_use_background_image"`
	// 		ProfileImageURL                string `json:"profile_image_url"`
	// 		ProfileImageURLHTTPS           string `json:"profile_image_url_https"`
	// 		DefaultProfile                 bool   `json:"default_profile"`
	// 		DefaultProfileImage            bool   `json:"default_profile_image"`
	// 		Following                      string `json:"following"`
	// 		FollowRequestSent              string `json:"follow_request_sent"`
	// 		Notifications                  string `json:"notifications"`
	// 	} `json:"user"`
	// 	Geo           string `json:"geo"`
	// 	Coordinates   string `json:"coordinates"`
	// 	Place         string `json:"place"`
	// 	Contributors  string `json:"contributors"`
	// 	IsQuoteStatus bool   `json:"is_quote_status"`
	// 	QuoteCount    int    `json:"quote_count"`
	// 	ReplyCount    int    `json:"reply_count"`
	// 	RetweetCount  int    `json:"retweet_count"`
	// 	FavoriteCount int    `json:"favorite_count"`
	// 	Entities      struct {
	// 		Hashtags     []interface{} `json:"hashtags"`
	// 		Urls         []interface{} `json:"urls"`
	// 		UserMentions []interface{} `json:"user_mentions"`
	// 		Symbols      []interface{} `json:"symbols"`
	// 	} `json:"entities"`
	// 	Favorited   bool   `json:"favorited"`
	// 	Retweeted   bool   `json:"retweeted"`
	// 	FilterLevel string `json:"filter_level"`
	// 	Lang        string `json:"lang"`
	// 	TimestampMs string `json:"timestamp_ms"`
	// } `json:"tweet"`
}

func NewFactomTweetFromContent(content []byte) (*FactomTweetContent, error) {
	tc := new(FactomTweetContent)
	// tc := new(twitter.Tweet)
	err := json.Unmarshal(content, tc)
	return tc, err
}

func (t *FactomTweetContent) TweetHash() []byte {
	hash := sha256.New().Sum([]byte(t.Tweet.Text))
	return hash
}

// "Thu Jan 24 15:01:51 +0000 2019"
func ParseTwitterDate(date string) (time.Time, error) {
	return time.Parse(time.RubyDate, date)
}
