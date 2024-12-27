package main

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
}
