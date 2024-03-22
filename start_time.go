package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	cmd "Bot_Compliment/command"
	event "Bot_Compliment/handleevent"
	send "Bot_Compliment/handlemessage"
	msg "Bot_Compliment/message"
	user "Bot_Compliment/userdata"
)

var (
	userData []user.UserData
)

func main() {
	bot, err := tgbotapi.NewBotAPI(getTikenFromArgument())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Активирован бот %s", bot.Self.UserName)

	user.LoadUserData(userData)
	event.LoadScheduledEvents()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Канал для хранения отложенных сообщений
	scheduledMessages := make(chan msg.ScheduledMessage)

	// Горутина для отправки отложенных сообщений
	go send.SendScheduledMessages(bot, scheduledMessages)

	for update := range updates {
		if update.Message == nil { // Игнорируем любые обновления, кроме сообщений
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				send.SendMessage(bot, update.Message.Chat.ID, "Добро пожаловать! Используйте команду /time для установки времени отправки сообщений.")
			case "time":
				go cmd.HandleTimeCommand(bot, update.Message, scheduledMessages, update.Message.Chat.UserName, userData)
			case "stop":
				go cmd.HandleStopCommand(bot, update.Message, userData)
			}
		}
	}
}

func getTikenFromArgument() string {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please give Token")
		panic("No arg")
	}
	return arguments[1]
}
