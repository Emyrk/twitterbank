package apiserver

import (
	"fmt"

	"github.com/Emyrk/twitterbank/database"
	"github.com/graphql-go/graphql"
)

var (
	APILimits = map[string][]int{
		"tweet_count": []int{50, 100},
	}
)

func getLimitAndOffset(api string, p graphql.ResolveParams) (int, int) {
	v := APILimits[api]
	def, max := v[0], v[1]
	o := 0
	if po, ok := p.Args["offset"]; ok {
		o = po.(int)
	}
	if l, ok := p.Args["limit"]; ok {
		if l.(int) > max {
			return max, o
		}
		return l.(int), o
	}
	return def, o
}

func (api *TwitterBankApiServer) TwitterUser() *graphql.Field {
	return &graphql.Field{
		Type:        api.TwitterUserType(),
		Description: "A user on Twitter.",
		Args: graphql.FieldConfigArgument{
			"user_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The uid of the twitter user",
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			uid := params.Args["user_id"].(string)
			return api.DB.FetchUserByUID(uid)
		},
	}
}

func (api *TwitterBankApiServer) TwitterUserType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "TwitterUser",
		Description: "Fetch a user on Twitter",
		Fields: graphql.Fields{
			"user_id_str": &graphql.Field{
				Type:        graphql.String,
				Description: "User ID of a given twitter user",
			},
			"tweets": &graphql.Field{
				Type:        graphql.NewList(TwitterTweetType),
				Description: "List all tweets by the given user.",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "If fetching tweets, use a limit to the number of tweets.",
					},
					"offset": &graphql.ArgumentConfig{
						Type:        graphql.Int,
						Description: "If fetching tweets, use a offset to the tweets retrieved.",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					tu, ok := p.Source.(*database.TwitterUser)
					if !ok {
						return nil, fmt.Errorf("Twitter user not found")
					}
					l, o := getLimitAndOffset("tweet_count", p)
					err := tu.FindTweets(api.DB.DB, l, o)
					return tu.Tweets, err
				},
			},
		}})
}

func (api *TwitterBankApiServer) Tweet() *graphql.Field {
	return &graphql.Field{
		Name:        "Tweet",
		Description: "Fetch a tweet on Twitter",
		Type:        TwitterTweetType,
		Args: graphql.FieldConfigArgument{
			"tweet_id": &graphql.ArgumentConfig{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "The tweet id of the wanted tweet.",
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			tid := p.Args["tweet_id"]
			return api.DB.FetchTweetByTID(tid.(string))
		},
	}
}

var TwitterTweetType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "TwitterTweet",
	Description: "A tweet on Twitter",
	Fields: graphql.Fields{
		"tweet_id_str": &graphql.Field{
			Type:        graphql.String,
			Description: "Tweet unique id.",
		},
		"tweet_author_str": &graphql.Field{
			Type:        graphql.String,
			Description: "User ID of the tweet author.",
		},
		"chain_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Chain ID of the author for this tweet.",
		},
		"entry_hash": &graphql.Field{
			Type:        graphql.String,
			Description: "Entryhash of the FIRST factom record of this tweet.",
		},
		"tweet_hash": &graphql.Field{
			Type:        graphql.String,
			Description: "SHA256 Hash of the tweet text",
		},
		"tweeted_time": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "Time the tweet was tweeted on the twitter platform.",
		},
	}})
