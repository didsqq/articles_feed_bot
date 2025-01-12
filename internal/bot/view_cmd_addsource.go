package bot

import (
	"context"
	"fmt"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	"github.com/didsqq/news_feed_bot/internal/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SourceStorage interface {
	Add(ctx context.Context, source model.Source) (int64, error)
}

func ViewCmdAddSource(storage SourceStorage) botkit.ViewFunc {
	type addSourceArgs struct {
		Name     string `json:"name"`
		URL      string `json:"url"`
		Priority int    `json:"priority"`
	}

	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		res, err := botkit.ParseJSON[addSourceArgs](update.Message.CommandArguments())
		if err != nil {
			return err
		}

		source := model.Source{
			Name:     res.Name,
			FeedURL:  res.URL,
			Priority: res.Priority,
		}

		sourceID, err := storage.Add(ctx, source)
		if err != nil {
			//send err msg
			return err
		}

		var (
			msgText = fmt.Sprintf(
				"Источник добавлен с ID: '%d'\\. Используйте этот ID для обновления источника или удаления\\.",
				sourceID,
			)
			reply = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
