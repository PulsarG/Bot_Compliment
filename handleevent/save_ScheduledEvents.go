package handleevent

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"Bot_Compliment/data"
)

func SaveScheduledEvents() {
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