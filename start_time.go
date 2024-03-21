package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	send "Bot_Compliment/handlemessage"
	msg "Bot_Compliment/message"
	user "Bot_Compliment/userdata"
	comm "Bot_Compliment/command"
)

// Структура для хранения отложенных сообщений
/* type ScheduledMessage struct {
	ChatID  int64
	Message string
} */

// Структура для хранения информации о пользователе
/* type UserData struct {
	Username string
	ChatID   int64
	Hour1    int
	Min1     int
	Hour2    int
	Min2     int
} */

var (
	defaultMessage  = "Ваше отложенное сообщение отправлено."
	userDataFile    = "userdata.json"
	scheduledFile   = "scheduled.json"
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

	loadUserData()
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
				go handleTimeCommand(bot, update.Message, scheduledMessages, update.Message.Chat.UserName)
			case "stop":
				go comm.HandleStopCommand(bot, update.Message, userData)
			}
		}
	}
}

// Функция для удаления пользователя и его данных
/* func handleStopCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
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
	saveUserData()

	// Отправка сообщения об удалении пользователя
	send.SendMessage(bot, chatID, "Вы успешно отписались от сервиса. Ваши данные удалены.")
} */

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

// Обработка команды /time для установки времени отправки сообщений
func handleTimeCommand(bot *tgbotapi.BotAPI, mssg *tgbotapi.Message, scheduledMessages chan msg.ScheduledMessage, username string) {
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
			scheduledMessages <- msg.ScheduledMessage{ChatID: mssg.Chat.ID, Message: defaultMessage}
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
				scheduledMessages <- msg.ScheduledMessage{ChatID: mssg.Chat.ID, Message: defaultMessage}
			}(duration, i)
		}
		user.SaveUserData(userData)
	} else {
		send.SendMessage(bot, mssg.Chat.ID, "Использование: /time <час1> <минута1> [<час2> <минута2>]")
	}
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

// Функция для отправки сообщения
/* func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
} */

// Функция для сохранения информации о пользователе в файл
/* func saveUserData() {
	jsonData, err := json.Marshal(userData)
	if err != nil {
		log.Println("Ошибка сериализации данных пользователя:", err)
		return
	}
	err = ioutil.WriteFile(userDataFile, jsonData, 0644)
	if err != nil {
		log.Println("Ошибка записи данных пользователя в файл:", err)
	}
} */

// Функция для загрузки информации о пользователях из файла
func loadUserData() {
	file, err := os.Open(userDataFile)
	if err != nil {
		log.Println("Ошибка открытия файла с данными пользователя:", err)
		return
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Ошибка чтения файла с данными пользователя:", err)
		return
	}

	err = json.Unmarshal(jsonData, &userData)
	if err != nil {
		log.Println("Ошибка десериализации данных пользователя:", err)
	}
}

// Функция для загрузки отложенных событий из файла
func loadScheduledEvents() {
	file, err := os.Open(scheduledFile)
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
	err = ioutil.WriteFile(scheduledFile, jsonData, 0644)
	if err != nil {
		log.Println("Ошибка записи отложенных событий в файл:", err)
	}
}
