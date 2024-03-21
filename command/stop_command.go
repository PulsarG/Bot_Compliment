package command

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"

	send "Bot_Compliment/handlemessage"
	user "Bot_Compliment/userdata"
)

func HandleStopCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, userData []user.UserData) {
	chatID := msg.Chat.ID
	username := msg.Chat.UserName

	// Удаление пользователя из данных
	for i, data := range userData {
		if data.Username == username {
			userData = append(userData[:i], userData[i+1:]...)
			break
		}
	}

	// Сохранение обновленных данных
	user.SaveUserData(userData)

	// Отправка сообщения об удалении пользователя
	send.SendMessage(bot, chatID, "Вы успешно отписались от сервиса. Ваши данные удалены.")
}
