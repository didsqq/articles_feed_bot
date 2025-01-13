package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	"github.com/didsqq/news_feed_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceStorageSources interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

func ViewCmdListSources(storage SourceStorageSources) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		sources, err := storage.Sources(ctx)
		if err != nil {
			return err
		}

		var msgText []string

		for _, source := range sources {
			msgText = append(msgText, fmt.Sprintf("Name:%s\nID:%d\nURL:%s\nPriority:%d\n\n\n", source.Name, source.ID, source.FeedURL, source.Priority))
		}
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(msgText, "\n"))

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
