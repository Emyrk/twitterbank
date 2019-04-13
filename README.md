# Twitter Seal

Note
```
docker-compose up -d --no-deps --build
```

# Factom structure

```bash
{
  "extids":[
    "TwitterBank Record", # Used to identify an entry for this project
    "TWITTER_HANDLE_ID", # To be able to find twitter user (use id not handle)
    "TWEET_ID", # Tweet id to locate tweet
    "IDENTITY_RECORDING", # Identity wotnessing tweet
    "IDENTITY_KEY",
    "SIGNATURE // Marshaled data excluding the sig (pad with 64 null bytes)",
  ],
  "content": {
    "dateFetched": "DATE_API_CALL",
    "tweet":{ # All data for the tweet that we want to keep
      "Tweet JSON",
    }
  }
}

```

You can use https://tweeterid.com/ to find a twitter handle id.

By including data such as the twitter handle id, and tweet id, it prevents replys across other chains

# Twitter Objects

- **Mentions**
- **Hashtags**
- Cashtags
- Media
- Native Media

# Questions:
- Can you edit a tweet? :: NO!
- Do we have any replay issues?
- Sybil issues?
- Is the tweet hash consistent?
- How to handle media? Images, videos?
- How to handle links? Wayback machine?
