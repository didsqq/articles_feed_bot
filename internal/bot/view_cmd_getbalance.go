package bot

import (
	"context"
	"fmt"

	"github.com/didsqq/news_feed_bot/internal/botkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BalanceGetProvider interface {
	GetBalance() (float64, error)
}

func ViewCmdGetBalance(api BalanceGetProvider) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		balance, err := api.GetBalance()
		if err != nil {
			return err
		}

		var (
			msgText = fmt.Sprintf("Баланс: %f", balance)
			reply   = tgbotapi.NewMessage(update.Message.Chat.ID, msgText)
		)

		if _, err := bot.Send(reply); err != nil {
			return err
		}

		return nil
	}
}
