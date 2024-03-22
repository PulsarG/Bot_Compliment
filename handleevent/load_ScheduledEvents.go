package handleevent

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"


	"Bot_Compliment/data"
)

func LoadScheduledEvents() {
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