package bot

import (
	"context"
	"fmt"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UsersGetStorage interface {
	GetKeywords(ctx context.Context, chatId int64) (string, error)
}

func ViewCmdGetKeywords(storage UsersGetStorage) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		keywords, err := storage.GetKeywords(ctx, update.Message.Chat.ID)
		if err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf("Ключевые слова: %s", keywords)
			reply   = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
