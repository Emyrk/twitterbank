package apiserver

import (
	"fmt"

	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"

	"github.com/Emyrk/twitterbank/database"
	"github.com/graphql-go/graphql"
)

var (
	APILimits = map[string][]int{
		"tweet_count": []int{50, 100},
	}
)

type GraphQLAPITypes struct {
	TwitterTweetSingleton        *graphql.Object
	TwitterUserSingletonNoTweets *graphql.Object
	TwitterUserSingleton         *graphql.Object
}

var TweetJSON = TweetJSONScalar()

func (api *TwitterBankApiServer) init() {
	api.apiTypes = new(GraphQLAPITypes)
	api.apiTypes.TwitterTweetSingleton = api.TwitterTweetType()
	api.apiTypes.TwitterUserSingletonNoTweets = api.TwitterUserType()
	api.apiTypes.TwitterUserSingleton = api.TwitterUserType()
	api.apiTypes.TwitterUserSingleton.PrivateName = "TwitteUserAndTweets"
	api.apiTypes.TwitterUserSingleton.AddFieldConfig("tweets", &graphql.Field{
		Type:        graphql.NewList(api.apiTypes.TwitterTweetSingleton),
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
			fmt.Printf("%v\n", p.Context)
			fmt.Printf("%v\n", p.Args)
			fmt.Printf("%v\n", p.Info)
			fmt.Printf("%v\n", p.Source)
			tu, ok := p.Source.(*database.TwitterUser)
			if !ok {
				return nil, fmt.Errorf("Twitter user not found")
			}
			l, o := getLimitAndOffset("tweet_count", p)
			err := tu.FindTweets(api.DB.DB, l, o)
			return tu.Tweets, err
		},
	})
}

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

func (api *TwitterBankApiServer) TwitterUsers() *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(api.apiTypes.TwitterUserSingletonNoTweets),
		Description: "All users being tracked",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return api.DB.FetchAllUsers()
		},
	}
}

func (api *TwitterBankApiServer) TwitterUser() *graphql.Field {
	return &graphql.Field{
		Type:        api.apiTypes.TwitterUserSingleton,
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
			"chain_id": &graphql.Field{
				Type:        graphql.String,
				Description: "Chainid tracking the user",
			},
			"entry_hash": &graphql.Field{
				Type:        graphql.String,
				Description: "First entry for a user chain",
			},
			"block_initiated": &graphql.Field{
				Type:        graphql.Int,
				Description: "What block the user tracking started in",
			},
		}})
}

func (api *TwitterBankApiServer) Tweet() *graphql.Field {
	return &graphql.Field{
		Name:        "Tweet",
		Description: "Fetch a tweet on Twitter",
		Type:        api.apiTypes.TwitterTweetSingleton,
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

func (api *TwitterBankApiServer) TwitterTweetType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
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
			"raw_tweet": &graphql.Field{
				Type:        TweetJSON,
				Description: "This field has not been dissected yet",
			},
			"proofs": &graphql.Field{
				Type:        graphql.NewList(TwitterTweetProofType),
				Description: "Recorded proofs into the Factom blockchain",
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
					tu, ok := p.Source.(*database.TwitterTweetObject)
					if !ok {
						tuBase, ok := p.Source.(database.TwitterTweetObject)
						if !ok {
							return nil, fmt.Errorf("Tweet not found")
						}
						tu = &tuBase
					}
					l, o := getLimitAndOffset("tweet_count", p)
					err := tu.FindProofs(api.DB.DB, l, o)
					return tu.Proofs, err
				},
			},
		}})
}

var TwitterTweetProofType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "FactomProof",
	Description: "Recorded proof of tweet by a Factom Identity.",
	Fields: graphql.Fields{
		"identity": &graphql.Field{
			Type:        graphql.String,
			Description: "Recording Identity.",
		},
		"tweet_author_str": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique id of the author of the tweet, given by Twitter.com.",
		},
		"tweet_id_str": &graphql.Field{
			Type:        graphql.String,
			Description: "Unique id of the tweet, given by Twitter.com.",
		},
		"tweet_hash": &graphql.Field{
			Type:        graphql.String,
			Description: "SHA256 hash of the tweet's text.",
		},
		"chain_id": &graphql.Field{
			Type:        graphql.String,
			Description: "Chain id of the entry of proof.",
		},
		"entry_hash": &graphql.Field{
			Type:        graphql.String,
			Description: "Entry hash of the entry of proof.",
		},
		"signing_key": &graphql.Field{
			Type:        graphql.String,
			Description: "Signing key used by Identity",
		},
		"signature": &graphql.Field{
			Type:        graphql.String,
			Description: "Signature of record",
		},
		"tweet_recorded_time": &graphql.Field{
			Type:        graphql.String,
			Description: "Time recorded into the Factom Blockchain",
		},
		"block_height": &graphql.Field{
			Type:        graphql.String,
			Description: "Block height of record in Factom Blockchain",
		},
	},
})

func parseLiteral(astValue ast.Value) interface{} {
	kind := astValue.GetKind()

	switch kind {
	case kinds.StringValue:
		return astValue.GetValue()
	case kinds.BooleanValue:
		return astValue.GetValue()
	case kinds.IntValue:
		return astValue.GetValue()
	case kinds.FloatValue:
		return astValue.GetValue()
	case kinds.ObjectValue:
		obj := make(map[string]interface{})
		for _, v := range astValue.GetValue().([]*ast.ObjectField) {
			obj[v.Name.Value] = parseLiteral(v.Value)
		}
		return obj
	case kinds.ListValue:
		list := make([]interface{}, 0)
		for _, v := range astValue.GetValue().([]ast.Value) {
			list = append(list, parseLiteral(v))
		}
		return list
	default:
		return nil
	}
}

// // JSON json type
// var JSON = graphql.NewScalar(
// 	graphql.ScalarConfig{
// 		Name:        "JSON",
// 		Description: "The `JSON` scalar type represents JSON values as specified by [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf)",
// 		Serialize: func(value interface{}) interface{} {
// 			return value
// 		},
// 		ParseValue: func(value interface{}) interface{} {
// 			return value
// 		},
// 		ParseLiteral: parseLiteral,
// 	},
// )

func TweetJSONScalar() *graphql.Scalar {
	return graphql.NewScalar(
		graphql.ScalarConfig{
			Name:        fmt.Sprintf("TweetJSON"),
			Description: "The `JSON` scalar type represents JSON values as specified by [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf)",
			Serialize: func(value interface{}) interface{} {
				s, ok := value.(string)
				if !ok {
					panic("sss")
					return nil
				}
				tu, err := database.NewFactomTweetFromContent([]byte(s))
				if err != nil {
					panic(err.Error())
					return nil
				}
				return tu
			},
			ParseValue: func(value interface{}) interface{} {
				return value
			},
			ParseLiteral: parseLiteral,
		},
	)
}
