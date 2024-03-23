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

	// !! Реализовать через формат HH:MM HH:MM
	// Разделение параметров по пробелам
	params := strings.Split(args, " ")

	if len(params) == 2 { // Пользователь ввел только одно время
		readParamsTimeCommandAndSave(1, params, &userData, bot, mssg.Chat.ID, username, scheduledMessages)
		user.SaveUserData(userData)
	} else if len(params) == 4 { // Пользователь ввел два времени
		readParamsTimeCommandAndSave(2, params, &userData, bot, mssg.Chat.ID, username, scheduledMessages)
		user.SaveUserData(userData)
	} else {
		send.SendMessage(bot, mssg.Chat.ID, data.Wrong_timeCommand)
	}
}

func readParamsTimeCommandAndSave(countParams int, params []string, userDataPoint *[]user.UserData, bot *tgbotapi.BotAPI, userID int64, username string, scheduledMessages chan msg.ScheduledMessage) {
	userData := *userDataPoint
	hours := make([]int, countParams)
	minutes := make([]int, countParams)
	for i := 0; i < countParams; i++ {
		hour, err := strconv.Atoi(params[i*2])
		if err != nil {
			send.SendMessage(bot, userID, "Неверный час")
			return
		}
		hours[i] = hour

		minute, err := strconv.Atoi(params[i*2+1])
		if err != nil {
			send.SendMessage(bot, userID, "Неверная минута")
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
	if !found && countParams == 2 {
		userData = append(userData, user.UserData{Username: username, ChatID: userID, Hour1: hours[0], Min1: minutes[0], Hour2: hours[1], Min2: minutes[1]})
	} else {
		userData = append(userData, user.UserData{Username: username, ChatID: userID, Hour1: hours[0], Min1: minutes[0]})
	}

	// Установка отложенных времен для отправки сообщений
	for i := 0; i < countParams; i++ {
		targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours[i], minutes[i], 0, 0, time.Local)
		duration := targetTime.Sub(time.Now())
		go func(d time.Duration, index int) {
			time.Sleep(d)
			scheduledMessages <- msg.ScheduledMessage{ChatID: userID, Message: data.DefaultMessage}
		}(duration, i)
	}

	*userDataPoint = userData
}
