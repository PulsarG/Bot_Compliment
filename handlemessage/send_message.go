package handlemessage

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func SendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
