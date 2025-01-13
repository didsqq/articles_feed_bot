package bot

import (
	"context"
	"fmt"
	"strconv"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceStorageDelete interface {
	Delete(ctx context.Context, id int64) error
}

func ViewCmdDeleteSource(storage SourceStorageDelete) botkit.ViewFunc {

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		idStr := update.Message.CommandArguments()

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return err
		}

		err = storage.Delete(ctx, id)
		if err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf(
				"Источник удален с ID: '%d'\\.",
				id,
			)
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
