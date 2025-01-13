package bot

import (
	"context"
	"errors"
	"strings"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	"github.com/didsqq/news_feed_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UsersAddKeyStorage interface {
	AddKeywords(ctx context.Context, user model.User) error
}

func ViewCmdAddKeywords(storage UsersAddKeyStorage) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		keyWords := update.Message.CommandArguments()

		if keyWords == "" {
			return errors.New("no keys in argument")
		}
		user := model.User{
			ChatID:   update.Message.Chat.ID,
			Keywords: strings.Split(keyWords, " "),
		}

		err := storage.AddKeywords(ctx, user)
		if err != nil {
			return err
		}

		var (
			msgText = "Ключевые слова успешно добавлены"
			reply   = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
