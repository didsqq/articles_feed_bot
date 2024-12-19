go install github.com/pressly/goose/v3/cmd/goose@latest
goose create init sql - создание миграции
goose postgres "host=localhost user=postgres database=news_feed_bot password=postgres sslmode=disable" status - проверка статуса
goose postgres "host=localhost user=postgres database=news_feed_bot password=postgres sslmode=disable" up - применение