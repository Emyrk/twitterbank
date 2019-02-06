# Twitter Seal

Each tweet object:
```json
{
  "created_at": "Thu Apr 06 15:24:15 +0000 2017",
  "id_str": "850006245121695744",
  "text": "1\/ Today we\u2019re sharing our vision for the future of the Twitter API platform!\nhttps:\/\/t.co\/XweGngmxlP",
  "user": {
    "id": 2244994945,
    "name": "Twitter Dev",
    "screen_name": "TwitterDev",
    "location": "Internet",
    "url": "https:\/\/dev.twitter.com\/",
    "description": "Your official source for Twitter Platform news, updates & events. Need technical help? Visit https:\/\/twittercommunity.com\/ \u2328\ufe0f #TapIntoTwitter"
  },
  "place": {   
  },
  "entities": {
    "hashtags": [      
    ],
    "urls": [
      {
        "url": "https:\/\/t.co\/XweGngmxlP",
        "unwound": {
          "url": "https:\/\/cards.twitter.com\/cards\/18ce53wgo4h\/3xo1c",
          "title": "Building the Future of the Twitter API Platform"
        }
      }
    ],
    "user_mentions": [     
    ]
  }
}
```

There are many other fields, such as `extended_tweet`, `retweeted_status`, etc.

The full tweet should always be recorded, especially as twitter expands their api with more meta data.
We should attempt to keep our structures independent of twitter's if we can. Especially for our back end.

Will we have to interpret quotes and retweets?


# Factom structure

```json
{
  "extids":[
    "TwitterBank Record",
    "TWITTER_HANDLE",
    "TWEET_HASH",
    "IDENTITY_RECORDING",
    "IDENTITY_KEY",
    "SIGNATURE // Marshaled data excluding the sig (pad with 64 null bytes)",
  ],
  "content": {
    "dateFetched": "DATE_API_CALL",
    "tweet":{
      "Tweet JSON",
    }
  }
}

```

# Questions:
- Can you edit a tweet? :: NO!
- Do we have any replay issues?
- Sybil issues?
- Is the tweet hash consistent?
- How to handle media? Images, videos?
- How to handle links? Wayback machine?