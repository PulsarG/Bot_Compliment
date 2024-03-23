package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	cmd "Bot_Compliment/command"
	"Bot_Compliment/data"
	send "Bot_Compliment/handlemessage"
	msg "Bot_Compliment/message"
	user "Bot_Compliment/userdata"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(getTokenFromArgument())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Активирован бот %s", bot.Self.UserName)

	user.LoadUserData(user.UserDatas)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

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
			case data.Command_Start:
				send.SendMessage(bot, update.Message.Chat.ID, data.Message_Welcome)
			case data.Command_Time:
				go cmd.HandleTimeCommand(bot, update.Message, scheduledMessages, update.Message.Chat.UserName, user.UserDatas)
			case data.Command_Stop:
				go cmd.HandleStopCommand(bot, update.Message, user.UserDatas)
			}
		}
	}
}

func getTokenFromArgument() string {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please give Token")
		panic("No arg")
	}
	return arguments[1]
}
