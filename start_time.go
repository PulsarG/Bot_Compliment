package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

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
	go sendScheduledMessages(bot, scheduledMessages)

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

func getRandomUniqueMessage(filePath string, chosenMessages map[string]bool) (string, error) {
	// Чтение содержимого файла
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Разделение содержимого файла на отдельные сообщения по переносу строки
	messages := strings.Split(string(content), "\n")

	// Создание списка доступных для выбора сообщений (которые еще не были выбраны)
	availableMessages := make([]string, 0)
	for _, message := range messages {
		if !chosenMessages[message] && message != "" {
			availableMessages = append(availableMessages, message)
		}
	}

	// Если доступных сообщений нет, возвращаем ошибку
	if len(availableMessages) == 0 {
		return "", errors.New("все сообщения уже были выбраны")
	}

	// Выбор рандомного сообщения из доступных
	rand.Seed(time.Now().UnixNano())
	selectedMessage := availableMessages[rand.Intn(len(availableMessages))]

	// Отметка выбранного сообщения как выбранного
	chosenMessages[selectedMessage] = true

	return selectedMessage, nil
}

// Горутина для отправки отложенных сообщений
func sendScheduledMessages(bot *tgbotapi.BotAPI, scheduledMessages chan msg.ScheduledMessage) {
	for {
		msg := <-scheduledMessages

		// !! TO DO Вынести в отдельную функцию

		var message string
		filePath := "messages.txt"
		chosenMessages := make(map[string]bool)
		for i := 0; i < 1; i++ {
			randMessage, err := getRandomUniqueMessage(filePath, chosenMessages)
			if err != nil {
				log.Println("Ошибка выбора сообщения:", err)
				break
			}
			message = randMessage
			log.Println("Выбранное сообщение:", message)
		}

		send.SendMessage(bot, msg.ChatID, message)
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
