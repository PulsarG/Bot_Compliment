package main

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type UserSchedule struct {
	ChatID    int64
	Schedule  []ScheduledMessage
}

type ScheduledMessage struct {
	Time     time.Time
	Message  string
}

var (
	scheduleMap  = make(map[int64]*UserSchedule)
	scheduleLock sync.Mutex
)

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if strings.HasPrefix(update.Message.Text, "/set") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите время (чч:мм) и сообщение через пробел, например, /set 10:30 Позвонить маме")
			bot.Send(msg)
		}

		if strings.HasPrefix(update.Message.Text, "/set ") {
			msgText := strings.SplitN(update.Message.Text, " ", 3)
			if len(msgText) != 3 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный формат. Используйте /set HH:MM Текст сообщения")
				bot.Send(msg)
				continue
			}

			t, err := time.Parse("15:04", msgText[1])
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный формат времени. Используйте HH:MM")
				bot.Send(msg)
				continue
			}

			userSchedule := &UserSchedule{
				ChatID: update.Message.Chat.ID,
			}

			scheduleLock.Lock()
			if existingSchedule, ok := scheduleMap[update.Message.Chat.ID]; ok {
				userSchedule = existingSchedule
			} else {
				scheduleMap[update.Message.Chat.ID] = userSchedule
			}
			scheduleLock.Unlock()

			userSchedule.Schedule = append(userSchedule.Schedule, ScheduledMessage{
				Time:    t,
				Message: msgText[2],
			})

			go func() {
				// Рассчитываем время до отправки сообщения
				now := time.Now()
				loc, _ := time.LoadLocation("Local")
				today := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, loc)
				if today.Before(now) {
					today = today.AddDate(0, 0, 1)
				}

				duration := today.Sub(now)

				// Ждем указанное время
				time.Sleep(duration)

				// Отправляем сообщение
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText[2])
				bot.Send(msg)

				// Удаляем сообщение из расписания
				scheduleLock.Lock()
				userSchedule.Schedule = userSchedule.Schedule[1:]
				if len(userSchedule.Schedule) == 0 {
					delete(scheduleMap, update.Message.Chat.ID)
				}
				scheduleLock.Unlock()
			}()
		}
	}
}
