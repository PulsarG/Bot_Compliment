package handleevent

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"

	"Bot_Compliment/data"
)

func SaveScheduledEvents() {
	fmt.Print("Start SAVE shed")
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