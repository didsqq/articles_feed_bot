# RSS Feed Bot
A bot for Telegram that receives articles via RSS feed makes summary using gpt and sends them to users using keywords and to the channel

Status: In development
## Features
- Fetching articles from RSS feeds
- Article summaries powered by gpt-4o-mini
- Admin commands for managing sources and api balance
- User commands for managing keys
## Configuration
### Environment variables
- `AFB_TELEGRAM_BOT_TOKEN` — token for Telegram Bot API
- `AFB_TELEGRAM_CHANNEL_ID` — ID of the channel to post to, can be obtained via @JsonDumpBot
- `AFB_DATABASE_DSN` — PostgreSQL connection string
- `AFB_FETCH_INTERVAL` — the interval of checking for new articles, default 10m
- `AFB_NOTIFICATION_INTERVAL` — the interval of delivering new articles to Telegram channel, default 1m
- `AFB_FILTER_KEYWORDS` — comma separated list of words to skip articles containing these words
- `AFB_OPENAI_KEY` — token for OpenAI API
- `AFB_OPENAI_PROMPT` — prompt for GPT to generate summary
### HCL
News Feed Bot can be configured with HCL config file. The service is looking for config file in following locations:

    ./config.hcl
    ./config.local.hcl
    $HOME/.config/news-feed-bot/config.hcl
## Technologies used
- Go
- PostgreSQL
- Docker
- OpenAI API
- Telegram Bot API
