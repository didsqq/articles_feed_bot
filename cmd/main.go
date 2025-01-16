package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/didsqq/news_feed_bot/internal/api"
	"github.com/didsqq/news_feed_bot/internal/bot"
	"github.com/didsqq/news_feed_bot/internal/bot/middleware"
	"github.com/didsqq/news_feed_bot/internal/botkit"
	"github.com/didsqq/news_feed_bot/internal/config"
	"github.com/didsqq/news_feed_bot/internal/fetcher"
	"github.com/didsqq/news_feed_bot/internal/notifier"
	"github.com/didsqq/news_feed_bot/internal/storage"
	"github.com/didsqq/news_feed_bot/internal/summary"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("failed to create bot: %v", err)
		return
	}

	db, err := sqlx.Connect("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return
	}
	defer db.Close()

	var (
		userStorage    = storage.NewUserStorage(db)
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		openAIClient   = api.NewOpenAIClient(
			config.Get().OpenAIKey,
			config.Get().OpenAIPrompt,
			config.Get().OpenAIModel,
		)
		articleFetcher = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)
		summarizer = summary.NewOpenAIProxySummarizer(
			openAIClient,
			config.Get().OpenAIModel,
		)
		articleNotifier = notifier.New(
			articleStorage,
			userStorage,
			summarizer,
			botAPI,
			config.Get().NotificationInterval,
			5*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	newsBot := botkit.New(botAPI)

	newsBot.RegisterCmdView(
		"balance",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdGetBalance(openAIClient),
		),
	)
	newsBot.RegisterCmdView(
		"addsource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdAddSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"listsource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdListSources(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"deletesource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdDeleteSource(sourceStorage),
		),
	)

	newsBot.RegisterCmdView("start", bot.ViewCmdStart(userStorage))
	newsBot.RegisterCmdView("addkeys", bot.ViewCmdAddKeywords(userStorage))
	newsBot.RegisterCmdView("getkeys", bot.ViewCmdGetKeywords(userStorage))
	newsBot.RegisterCmdView("deletekeys", bot.ViewCmdDeleteKeywords(userStorage))
	newsBot.RegisterCmdView("delete", bot.ViewCmdDeleteUser(userStorage))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func(ctx context.Context) {
		if err := articleFetcher.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to start fetcher: %v", err)
				return
			}
			log.Printf("[INFO] fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := articleNotifier.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to run notifier: %v", err)
				return
			}
			log.Printf("[INFO] notifier stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] failed to run botkit: %v", err)
		}
		log.Printf("[INFO] bot stopped")
	}
}
