package bot

import (
	"context"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UsersDeleteStorage interface {
	Delete(ctx context.Context, chatId int64) error
}

func ViewCmdDeleteUser(storage UsersDeleteStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {

		err := storage.Delete(ctx, update.Message.Chat.ID)
		if err != nil {
			return err
		}

		var (
			msgText = "Пользователь удален"
			reply   = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
