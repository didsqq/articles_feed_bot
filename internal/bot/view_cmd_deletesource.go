package bot

import (
	"context"
	"fmt"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceStorageDelete interface {
	Delete(ctx context.Context, id int64) error
}

func ViewCmdDeleteSource(storage SourceStorageDelete) botkit.ViewFunc {
	type deleteSourceArgs struct {
		Id int64 `json:"id"`
	}

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		res, err := botkit.ParseJSON[deleteSourceArgs](update.Message.CommandArguments())
		if err != nil {
			return err
		}

		err = storage.Delete(ctx, res.Id)
		if err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf(
				"Источник удален с ID: '%d'\\.",
				res.Id,
			)
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
