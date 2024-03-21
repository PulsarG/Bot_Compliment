package message

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

func GetRandomUniqueMessage(filePath string, chosenMessages map[string]bool) (string, error) {
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
