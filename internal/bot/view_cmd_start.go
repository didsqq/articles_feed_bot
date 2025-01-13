package bot

import (
	"context"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	"github.com/didsqq/news_feed_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UsersAddStorage interface {
	Add(ctx context.Context, user model.User) error
}

func ViewCmdStart(storage UsersAddStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		user := model.User{
			ChatID: update.Message.Chat.ID,
		}

		err := storage.Add(ctx, user)
		if err != nil {
			return err
		}

		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, "Hello, world!")); err != nil {
			return err
		}
		return nil
	}
}
