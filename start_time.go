package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	comm "Bot_Compliment/command"
	"Bot_Compliment/data"
	send "Bot_Compliment/handlemessage"
	msg "Bot_Compliment/message"
	user "Bot_Compliment/userdata"
)

var (
	userData        []user.UserData
	scheduledEvents []msg.ScheduledMessage
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please give Token")
		return
	}

	token := arguments[1]

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	//bot.Debug = true

	log.Printf("Активирован бот %s", bot.Self.UserName)

	user.LoadUserData(userData)
	loadScheduledEvents()

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
				go comm.HandleTimeCommand(bot, update.Message, scheduledMessages, update.Message.Chat.UserName, userData)
			case "stop":
				go comm.HandleStopCommand(bot, update.Message, userData)
			}
		}
	}
}

// Функция для загрузки отложенных событий из файла
func loadScheduledEvents() {
	file, err := os.Open(data.ScheduledFile)
	if err != nil {
		log.Println("Ошибка открытия файла с отложенными событиями:", err)
		return
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Ошибка чтения файла с отложенными событиями:", err)
		return
	}

	err = json.Unmarshal(jsonData, &scheduledEvents)
	if err != nil {
		log.Println("Ошибка десериализации отложенных событий:", err)
	}
}

// Функция для сохранения отложенных событий в файл
func saveScheduledEvents() {
	jsonData, err := json.Marshal(scheduledEvents)
	if err != nil {
		log.Println("Ошибка сериализации отложенных событий:", err)
		return
	}
	err = ioutil.WriteFile(data.ScheduledFile, jsonData, 0644)
	if err != nil {
		log.Println("Ошибка записи отложенных событий в файл:", err)
	}
}
