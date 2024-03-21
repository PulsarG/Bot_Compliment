package handlemessage

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	msg "Bot_Compliment/message"
)

func SendScheduledMessages(bot *tgbotapi.BotAPI, scheduledMessages chan msg.ScheduledMessage) {
	for {
		mssg := <-scheduledMessages

		// !! TO DO Вынести в отдельную функцию

		var message string
		filePath := "messages.txt"
		chosenMessages := make(map[string]bool)
		for i := 0; i < 1; i++ {
			randMessage, err := msg.GetRandomUniqueMessage(filePath, chosenMessages)
			if err != nil {
				log.Println("Ошибка выбора сообщения:", err)
				break
			}
			message = randMessage
			log.Println("Выбранное сообщение:", message)
		}

		SendMessage(bot, mssg.ChatID, message)
	}
}
