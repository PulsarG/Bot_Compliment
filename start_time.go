/* // 5708011095:AAHJiuyPCem8MSmZqbKpJCFzR11xT3lEwIk

package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Структура для хранения отложенных сообщений
type ScheduledMessage struct {
	ChatID  int64
	Message string
}

func main() {
	bot, err := tgbotapi.NewBotAPI("5708011095:AAHJiuyPCem8MSmZqbKpJCFzR11xT3lEwIk")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Авторизован на аккаунте %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Канал для хранения отложенных сообщений
	scheduledMessages := make(chan ScheduledMessage)

	// Горутина для отправки отложенных сообщений
	go sendScheduledMessages(bot, scheduledMessages)

	for update := range updates {
		if update.Message == nil { // Игнорируем любые обновления, кроме сообщений
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				sendMessage(bot, update.Message.Chat.ID, "Добро пожаловать! Используйте команду /time для установки времени отправки сообщений.")
			case "time":
				go handleTimeCommand(bot, update.Message, scheduledMessages)
			}
		}
	}
}

// Обработка команды /time для установки времени отправки сообщений
func handleTimeCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, scheduledMessages chan ScheduledMessage) {
	// Парсинг параметров из команды
	args := msg.CommandArguments()

	// Разделение параметров по пробелам
	params := strings.Split(args, " ")

	if len(params) == 2 { // Пользователь ввел только одно время
		hour, err := strconv.Atoi(params[0])
		if err != nil {
			sendMessage(bot, msg.Chat.ID, "Неверный час")
			return
		}

		minute, err := strconv.Atoi(params[1])
		if err != nil {
			sendMessage(bot, msg.Chat.ID, "Неверная минута")
			return
		}

		// Установка отложенного времени для отправки сообщения
		targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
		duration := targetTime.Sub(time.Now())
		go func(d time.Duration) {
			time.Sleep(d)
			scheduledMessages <- ScheduledMessage{ChatID: msg.Chat.ID, Message: "Ваше отложенное сообщение отправлено."}
		}(duration)
	} else if len(params) == 4 { // Пользователь ввел два времени
		// Разбор часов и минут из параметров
		hours := make([]int, 2)
		minutes := make([]int, 2)
		for i := 0; i < 2; i++ {
			hour, err := strconv.Atoi(params[i*2])
			if err != nil {
				sendMessage(bot, msg.Chat.ID, "Неверный час")
				return
			}
			hours[i] = hour

			minute, err := strconv.Atoi(params[i*2+1])
			if err != nil {
				sendMessage(bot, msg.Chat.ID, "Неверная минута")
				return
			}
			minutes[i] = minute
		}

		// Установка отложенных времен для отправки сообщений
		for i := 0; i < 2; i++ {
			targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours[i], minutes[i], 0, 0, time.Local)
			duration := targetTime.Sub(time.Now())
			go func(d time.Duration, index int) {
				time.Sleep(d)
				scheduledMessages <- ScheduledMessage{ChatID: msg.Chat.ID, Message: "Ваше отложенное сообщение " + strconv.Itoa(index+1) + " отправлено."}
			}(duration, i)
		}
	} else {
		sendMessage(bot, msg.Chat.ID, "Использование: /time <час1> <минута1> [<час2> <минута2>]")
	}
}

// Горутина для отправки отложенных сообщений
func sendScheduledMessages(bot *tgbotapi.BotAPI, scheduledMessages chan ScheduledMessage) {
	for {
		msg := <-scheduledMessages
		sendMessage(bot, msg.ChatID, msg.Message)
	}
}

// Функция для отправки сообщения
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
*/

/* package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Структура для хранения отложенных сообщений
type ScheduledMessage struct {
	ChatID  int64
	Message string
}

var defaultMessage = "Ваше отложенное сообщение отправлено."

func main() {
	bot, err := tgbotapi.NewBotAPI("5708011095:AAHJiuyPCem8MSmZqbKpJCFzR11xT3lEwIk")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Авторизован на аккаунте %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Канал для хранения отложенных сообщений
	scheduledMessages := make(chan ScheduledMessage)

	// Горутина для отправки отложенных сообщений
	go sendScheduledMessages(bot, scheduledMessages)

	for update := range updates {
		if update.Message == nil { // Игнорируем любые обновления, кроме сообщений
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				sendMessage(bot, update.Message.Chat.ID, "Добро пожаловать! Используйте команду /time для установки времени отправки сообщений.")
			case "time":
				go handleTimeCommand(bot, update.Message, scheduledMessages)
			case "setmessage":
				go handleSetMessageCommand(bot, update.Message)
			}
		}
	}
}

// Обработка команды /time для установки времени отправки сообщений
func handleTimeCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, scheduledMessages chan ScheduledMessage) {
	// Парсинг параметров из команды
	args := msg.CommandArguments()

	// Разделение параметров по пробелам
	params := strings.Split(args, " ")

	if len(params) == 2 { // Пользователь ввел только одно время
		hour, err := strconv.Atoi(params[0])
		if err != nil {
			sendMessage(bot, msg.Chat.ID, "Неверный час")
			return
		}

		minute, err := strconv.Atoi(params[1])
		if err != nil {
			sendMessage(bot, msg.Chat.ID, "Неверная минута")
			return
		}

		// Установка отложенного времени для отправки сообщения
		targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
		duration := targetTime.Sub(time.Now())
		go func(d time.Duration) {
			time.Sleep(d)
			scheduledMessages <- ScheduledMessage{ChatID: msg.Chat.ID, Message: defaultMessage}
		}(duration)
	} else if len(params) == 4 { // Пользователь ввел два времени
		// Разбор часов и минут из параметров
		hours := make([]int, 2)
		minutes := make([]int, 2)
		for i := 0; i < 2; i++ {
			hour, err := strconv.Atoi(params[i*2])
			if err != nil {
				sendMessage(bot, msg.Chat.ID, "Неверный час")
				return
			}
			hours[i] = hour

			minute, err := strconv.Atoi(params[i*2+1])
			if err != nil {
				sendMessage(bot, msg.Chat.ID, "Неверная минута")
				return
			}
			minutes[i] = minute
		}

		// Установка отложенных времен для отправки сообщений
		for i := 0; i < 2; i++ {
			targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours[i], minutes[i], 0, 0, time.Local)
			duration := targetTime.Sub(time.Now())
			go func(d time.Duration, index int) {
				time.Sleep(d)
				scheduledMessages <- ScheduledMessage{ChatID: msg.Chat.ID, Message: defaultMessage}
			}(duration, i)
		}
	} else {
		sendMessage(bot, msg.Chat.ID, "Использование: /time <час1> <минута1> [<час2> <минута2>]")
	}
}

// Обработка команды /setmessage для установки нового сообщения
func handleSetMessageCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMessage := msg.CommandArguments()
	if newMessage == "" {
		sendMessage(bot, msg.Chat.ID, "Использование: /setmessage <новое_сообщение>")
		return
	}
	defaultMessage = newMessage
	sendMessage(bot, msg.Chat.ID, "Сообщение успешно изменено")
}

// Горутина для отправки отложенных сообщений
func sendScheduledMessages(bot *tgbotapi.BotAPI, scheduledMessages chan ScheduledMessage) {
	for {
		msg := <-scheduledMessages
		sendMessage(bot, msg.ChatID, msg.Message)
	}
}

// Функция для отправки сообщения
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}
*/

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Структура для хранения отложенных сообщений
type ScheduledMessage struct {
	ChatID  int64
	Message string
}

// Структура для хранения информации о пользователе
type UserData struct {
	Username string
	ChatID   int64
	Hour1    int
	Min1     int
	Hour2    int
	Min2     int
}

var (
	defaultMessage  = "Ваше отложенное сообщение отправлено."
	userDataFile    = "userdata.json"
	scheduledFile   = "scheduled.json"
	userData        []UserData
	scheduledEvents []ScheduledMessage
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5708011095:AAHJiuyPCem8MSmZqbKpJCFzR11xT3lEwIk")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Авторизован на аккаунте %s", bot.Self.UserName)

	loadUserData()
	loadScheduledEvents()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	// Канал для хранения отложенных сообщений
	scheduledMessages := make(chan ScheduledMessage)

	// Горутина для отправки отложенных сообщений
	go sendScheduledMessages(bot, scheduledMessages)

	for update := range updates {
		if update.Message == nil { // Игнорируем любые обновления, кроме сообщений
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				sendMessage(bot, update.Message.Chat.ID, "Добро пожаловать! Используйте команду /time для установки времени отправки сообщений.")
			case "time":
				go handleTimeCommand(bot, update.Message, scheduledMessages, update.Message.Chat.UserName)
			case "setmessage":
				go handleSetMessageCommand(bot, update.Message)
			}
		}
	}
}

// Обработка команды /time для установки времени отправки сообщений
func handleTimeCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, scheduledMessages chan ScheduledMessage, username string) {
	// Парсинг параметров из команды
	args := msg.CommandArguments()

	// Разделение параметров по пробелам
	params := strings.Split(args, " ")

	if len(params) == 2 { // Пользователь ввел только одно время
		hour, err := strconv.Atoi(params[0])
		if err != nil {
			sendMessage(bot, msg.Chat.ID, "Неверный час")
			return
		}

		minute, err := strconv.Atoi(params[1])
		if err != nil {
			sendMessage(bot, msg.Chat.ID, "Неверная минута")
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
			userData = append(userData, UserData{Username: username, ChatID: msg.Chat.ID, Hour1: hour, Min1: minute})
		}

		// Установка отложенного времени для отправки сообщения
		targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hour, minute, 0, 0, time.Local)
		duration := targetTime.Sub(time.Now())
		go func(d time.Duration) {
			time.Sleep(d)
			scheduledMessages <- ScheduledMessage{ChatID: msg.Chat.ID, Message: defaultMessage}
		}(duration)
		saveUserData()
	} else if len(params) == 4 { // Пользователь ввел два времени
		// Разбор часов и минут из параметров
		hours := make([]int, 2)
		minutes := make([]int, 2)
		for i := 0; i < 2; i++ {
			hour, err := strconv.Atoi(params[i*2])
			if err != nil {
				sendMessage(bot, msg.Chat.ID, "Неверный час")
				return
			}
			hours[i] = hour

			minute, err := strconv.Atoi(params[i*2+1])
			if err != nil {
				sendMessage(bot, msg.Chat.ID, "Неверная минута")
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
			userData = append(userData, UserData{Username: username, ChatID: msg.Chat.ID, Hour1: hours[0], Min1: minutes[0], Hour2: hours[1], Min2: minutes[1]})
		}

		// Установка отложенных времен для отправки сообщений
		for i := 0; i < 2; i++ {
			targetTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours[i], minutes[i], 0, 0, time.Local)
			duration := targetTime.Sub(time.Now())
			go func(d time.Duration, index int) {
				time.Sleep(d)
				scheduledMessages <- ScheduledMessage{ChatID: msg.Chat.ID, Message: defaultMessage}
			}(duration, i)
		}
		saveUserData()
	} else {
		sendMessage(bot, msg.Chat.ID, "Использование: /time <час1> <минута1> [<час2> <минута2>]")
	}
}

// Обработка команды /setmessage для установки нового сообщения
func handleSetMessageCommand(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	newMessage := msg.CommandArguments()
	if newMessage == "" {
		sendMessage(bot, msg.Chat.ID, "Использование: /setmessage <новое_сообщение>")
		return
	}
	defaultMessage = newMessage
	sendMessage(bot, msg.Chat.ID, "Сообщение успешно изменено")
}

// Горутина для отправки отложенных сообщений
func sendScheduledMessages(bot *tgbotapi.BotAPI, scheduledMessages chan ScheduledMessage) {
	for {
		msg := <-scheduledMessages
		sendMessage(bot, msg.ChatID, msg.Message)
	}
}

// Функция для отправки сообщения
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("Ошибка отправки сообщения:", err)
	}
}

// Функция для сохранения информации о пользователе в файл
func saveUserData() {
	jsonData, err := json.Marshal(userData)
	if err != nil {
		log.Println("Ошибка сериализации данных пользователя:", err)
		return
	}
	err = ioutil.WriteFile(userDataFile, jsonData, 0644)
	if err != nil {
		log.Println("Ошибка записи данных пользователя в файл:", err)
	}
}

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
