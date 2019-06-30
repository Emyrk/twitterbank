package scraper

import (
	"strconv"

	"fmt"

	"github.com/Emyrk/twitterbank/database"
	"github.com/FactomProject/factomd/common/interfaces"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

var TestingString = "testing_kafka28"

var processLog = log.WithFields(log.Fields{"file": "processentry"})

type Processor struct {
	*Scraper // Need access to db and factom
}

func NewProcessor(s *Scraper) *Processor {
	p := new(Processor)
	p.Scraper = s
	return p
}

// ProcessEntry
//	Params:
//		entry		Contains data
//		dblock		Used to grab factom timing
func (s *Processor) ProcessEntry(entry interfaces.IEBEntry, dblock interfaces.IDirectoryBlock) error {
	if len(entry.ExternalIDs()) == 0 {
		return nil
	}
	switch string(entry.ExternalIDs()[0]) {
	case "TwitterBank Chain":
		// TODO: Remove this hack
		if len(entry.ExternalIDs()) == 3 {
			return s.ProcessTwitterChain(entry, dblock)
		}
		return s.ProcessTwitterEntry(entry, dblock)
	case "TwitterBank Record":
		return s.ProcessTwitterChain(entry, dblock)
	}

	return nil
}

func (p *Processor) ProcessTwitterChain(entry interfaces.IEBEntry, dblock interfaces.IDirectoryBlock) error {
	flog := processLog.WithFields(log.Fields{"func": "ProcessTwitterChain"})
	// Improper start to chain
	if len(entry.ExternalIDs()) != 3 {
		flog.Warnf("Chain %s has improper length extids to start chain", entry.GetChainID().String())
		return nil
	}

	if string(entry.ExternalIDs()[2]) != TestingString {
		return nil
	}

	//	"extids":[
	//		"TwitterBank Chain", # Used to identify an entry for this project
	//		"TWITTER_HANDLE_ID", # To be able to find twitter user (use id not handle)
	//	],

	handle_id := string(entry.ExternalIDs()[1])
	user_id, err := strconv.ParseInt(handle_id, 10, 64)
	if err != nil {
		// We don't really need this key.
		flog.Warnf("User_id to int failed: %s", err.Error())
	}

	user := database.TwitterUser{UserID: user_id, UserIDStr: handle_id, BlockInitiated: int(dblock.GetDatabaseHeight()), BlockInitatedUnix: dblock.GetTimestamp().GetTime().Unix(), ChainID: entry.GetChainID().String(), EntryHash: entry.GetHash().String()}

	return p.Database.InsertNewUserChain(&user)
}

const (
	RecordEntryTitle = iota
	RecordTwitterHandleID
	RecordTwitterTweetID
	// RecordIdentityRecording // TODO: Renable
	RecordIdentityKey
	RecordSignature
	RecordEnd
)

func (p *Processor) ProcessTwitterEntry(entry interfaces.IEBEntry, dblock interfaces.IDirectoryBlock) error {
	flog := processLog.WithFields(log.Fields{"func": "ProcessTwitterChain", "entry": entry.GetHash().String()})
	// Improper start to chain
	if len(entry.ExternalIDs()) != RecordEnd {
		flog.Warnf("Chain %s has improper length extids to start chain", entry.GetChainID().String())
		return nil
	}

	//{
	// 	"extids":[
	// 	0 "TwitterBank Record", # Used to identify an entry for this project
	// 	1 "TWITTER_HANDLE_ID", # To be able to find twitter user (use id not handle)
	// 	2 "TWEET_ID", # Tweet id to locate tweet
	// 	3 "IDENTITY_RECORDING", # Identity witnessing tweet
	// 	4 "IDENTITY_KEY",
	// 	5 "SIGNATURE // Marshaled data excluding the sig (pad with 64 null bytes)",
	// ],
	//		"content": {
	//		"dateFetched": "DATE_API_CALL",
	//		"tweet":{ # All data for the tweet that we want to keep
	//			"Tweet JSON",
	//		}
	//	}
	//}

	handle_id_str := string(entry.ExternalIDs()[RecordTwitterHandleID])
	handle_id, err := strconv.ParseInt(handle_id_str, 10, 64)
	if err != nil {
		// We don't really need this key.
		flog.Warnf("handle_id to int failed: %s", err.Error())
	}

	tweet_id_str := string(entry.ExternalIDs()[RecordTwitterTweetID])
	tweet_id, err := strconv.ParseInt(tweet_id_str, 10, 64)
	if err != nil {
		// We don't really need this key.
		flog.Warnf("tweet_id to int failed: %s", err.Error())
	}

	//TweetAuthorID    int64  `json:"tweet_author"`
	//TweetAuthorIDStr string `json:"tweet_author_str"`
	//
	//TweetID    int64  `json:"tweet_id"`
	//TweetIDStr string `json:"tweet_id_str" gorm:"primary_key"`
	//TweetHash  string `json:"tweet_hash" gorm:"primary_key"`
	//RawTweet string `json:"raw_tweet"`
	//TweetCreatedAt time.Time `json:"tweet_created_time"`

	tweet_content, err := database.NewFactomTweetFromContent(entry.GetContent())
	if err != nil {
		flog.Warnf("tweet_content to int failed: %s", err.Error())
		return nil
	}

	created_date, err := database.ParseTwitterDate(tweet_content.Tweet.CreatedAt)
	if err != nil {
		flog.Warnf("Twitter Date (%s) failed to parse: %s", tweet_content.Tweet.CreatedAt, err.Error())
		return nil
	}

	// TODO: Verify handle is in right chain
	tweet := database.TwitterTweetObject{
		TweetAuthorIDStr: handle_id_str,
		TweetAuthorID:    handle_id,
		TweetIDStr:       tweet_id_str,
		TweetID:          tweet_id,
		ChainID:          entry.GetChainID().String(),
		EntryHash:        entry.GetHash().String(),
		TweetCreatedAt:   created_date,
		RawTweet:         string(entry.GetContent()),
		Proofs:           []database.TwitterTweetRecord{},
	}

	// Some verification the tweet content matches the entry's header
	if !tweet.Collaborate(tweet_content) {
		flog.Warnf("Content does not collaborate extids Uid: %s/%s, Tid: %s/%s", tweet.TweetAuthorIDStr, handle_id_str, tweet.TweetIDStr, tweet_id_str)
		return nil
	}

	// TODO: Verify Identity
	// identity := fmt.Sprintf("%x", entry.ExternalIDs()[RecordIdentityRecording])
	identity := "placeholder"
	// TODO: Verify Signature
	// TODO: Verify key

	record := database.TwitterTweetRecord{
		FactomRecorder:   identity,
		EntryHash:        entry.GetHash().String(),
		ChainID:          entry.GetChainID().String(),
		TweetAuthorIDStr: handle_id_str,
		TweetAuthorID:    handle_id,
		TweetID:          tweet_id,
		TweetIDStr:       tweet_id_str,
		TweetHash:        string(tweet_content.TweetHash()),
		SigningKey:       fmt.Sprintf("%s", string(entry.ExternalIDs()[RecordIdentityKey])),
		Signature:        fmt.Sprintf("%s", string(entry.ExternalIDs()[RecordSignature])),
		TweetRecordedAt:  dblock.GetTimestamp().GetTime(),
		BlockHeight:      int(dblock.GetDatabaseHeight()),
	}

	err = p.Database.InsertNewTweet(&tweet, &record)
	if err != nil {
		return err
	}
	if tweet.Collaborate(tweet_content) {
		return err
	}

	return nil
}
