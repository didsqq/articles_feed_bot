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

		msg := `
Временно_для_админа
/addsource <ссылка на RSS> — Добавить новый источник статей (например, RSS-канал).
/listsource — Показать список всех добавленных источников.
/deletesource <номер> — Удалить источник из списка по его номеру.
/balance - Баланс gpt api.

Привет! 👋
Я бот, который поможет тебе получать самые интересные статьи из Хабра по ключевым словам и источникам!

Вот что я умею:

/start — Начать работу с ботом и получить это приветственное сообщение.
/addkeys <ключевые слова> — Добавить ключевые слова для фильтрации статей (вводите через запятую).
/getkeys — Посмотреть текущий список ключевых слов.
/deletekeys - Удалить ключевые слова
/delete — Остановить рассылку.
💡 Просто добавь источники и ключевые слова — и я начну присылать статьи, которые тебе подходят.
Если что-то непонятно, просто напиши мне, и я помогу! 😊`
		if _, err := bot.Send(tgbotapi.NewMessage(update.FromChat().ID, msg)); err != nil {
			return err
		}
		return nil
	}
}
