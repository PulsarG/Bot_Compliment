/* package main

import (
	"log"
	//"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Создание нового бота с токеном
	bot, err := tgbotapi.NewBotAPI("token")
	if err != nil {
		log.Fatal(err)
	}

	// Установка режима отладки (выводит запросы к API в консоль)
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Получение обновлений от бота
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Ответ на полученное сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		bot.Send(msg)
	}
}
 */