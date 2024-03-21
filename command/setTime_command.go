package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"Bot_Compliment/data"
	send "Bot_Compliment/handlemessage"
	msg "Bot_Compliment/message"
	user "Bot_Compliment/userdata"
)

func HandleTimeCommand(bot *tgbotapi.BotAPI, mssg *tgbotapi.Message, scheduledMessages chan msg.ScheduledMessage, username string, userData []user.UserData) {
	// Парсинг параметров из команды
	args := mssg.CommandArguments()

	// Разделение параметров по пробелам
	params := strings.Split(args, " ")

	if len(params) == 2 { // Пользователь ввел только одно время
		hour, err := strconv.Atoi(params[0])
		if err != nil {
			send.SendMessage(bot, mssg.Chat.ID, "Неверный час")
			return
		}

		minute, err := strconv.Atoi(params[1])
		if err != nil {
			send.SendMessage(bot, mssg.Chat.ID, "Неверная минута")
			return
		}

		// Поиск существующей записи для пользователя
		var found bool
		for i, data := range userData {
			if data.Username == username {
				userData[i].Hour1 = hour
				userData[i].Min1 = minute
				found = true
				break
			}
		}

		// Если запись не найдена, создаем новую
		if !found {
			userData = append(userData, user.UserData{Username: username, ChatID: mssg.Chat.ID, Hour1: hour, Min1: minute})
		}

		// Установка отложенного времени для отправки сообщения
		targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
		duration := targetTime.Sub(time.Now())
		go func(d time.Duration) {
			time.Sleep(d)
			scheduledMessages <- msg.ScheduledMessage{ChatID: mssg.Chat.ID, Message: data.DefaultMessage}
		}(duration)
		user.SaveUserData(userData)
	} else if len(params) == 4 { // Пользователь ввел два времени
		// Разбор часов и минут из параметров
		hours := make([]int, 2)
		minutes := make([]int, 2)
		for i := 0; i < 2; i++ {
			hour, err := strconv.Atoi(params[i*2])
			if err != nil {
				send.SendMessage(bot, mssg.Chat.ID, "Неверный час")
				return
			}
			hours[i] = hour

			minute, err := strconv.Atoi(params[i*2+1])
			if err != nil {
				send.SendMessage(bot, mssg.Chat.ID, "Неверная минута")
				return
			}
			minutes[i] = minute
		}

		// Поиск существующей записи для пользователя
		var found bool
		for i, data := range userData {
			if data.Username == username {
				userData[i].Hour1 = hours[0]
				userData[i].Min1 = minutes[0]
				userData[i].Hour2 = hours[1]
				userData[i].Min2 = minutes[1]
				found = true
				break
			}
		}

		// Если запись не найдена, создаем новую
		if !found {
			userData = append(userData, user.UserData{Username: username, ChatID: mssg.Chat.ID, Hour1: hours[0], Min1: minutes[0], Hour2: hours[1], Min2: minutes[1]})
		}

		// Установка отложенных времен для отправки сообщений
		for i := 0; i < 2; i++ {
			targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours[i], minutes[i], 0, 0, time.Local)
			duration := targetTime.Sub(time.Now())
			go func(d time.Duration, index int) {
				time.Sleep(d)
				scheduledMessages <- msg.ScheduledMessage{ChatID: mssg.Chat.ID, Message: data.DefaultMessage}
			}(duration, i)
		}
		user.SaveUserData(userData)
	} else {
		send.SendMessage(bot, mssg.Chat.ID, "Использование: /time <час1> <минута1> [<час2> <минута2>]")
	}
}