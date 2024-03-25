package userdata

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"Bot_Compliment/data"
)

func LoadUserData(userData *[]UserData) {
	file, err := os.Open(data.USER_DATA_FILE)
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
