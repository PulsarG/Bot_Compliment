package userdata

import (
	"encoding/json"
	"io/ioutil"
	"log"

	//"github.com/go-telegram-bot-api/telegram-bot-api"

	"Bot_Compliment/data"
	//user "Bot_Compliment/userdata"
)

func SaveUserData(userData []UserData) {
	jsonData, err := json.Marshal(userData)
	if err != nil {
		log.Println("Ошибка сериализации данных пользователя:", err)
		return
	}
	err = ioutil.WriteFile(data.UserDataFile, jsonData, 0644)
	if err != nil {
		log.Println("Ошибка записи данных пользователя в файл:", err)
	}
}
