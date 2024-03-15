package main

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// Создание нового бота с токеном
	bot, err := tgbotapi.NewBotAPI("5708011095:AAHJiuyPCem8MSmZqbKpJCFzR11xT3lEwIk")
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

	// Переменная для хранения статуса начала общения
	chatStarted := make(map[int64]bool)

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Обработка команды старт
		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я готов отвечать на ваши сообщения.")
			bot.Send(msg)

			// Устанавливаем статус начала общения для данного чата
			chatStarted[update.Message.Chat.ID] = true
			continue
		}

		// Проверяем, началось ли общение в данном чате
		if !chatStarted[update.Message.Chat.ID] {
			continue // Если нет, пропускаем обработку сообщения
		}

		// Ответ на полученное сообщение
		reply := "Вы сказали: " + update.Message.Text
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		bot.Send(msg)
	}
}
