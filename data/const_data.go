package data

import (
// user "Bot_Compliment/userdata"
)

var (
	DefaultMessage = "Ваше отложенное сообщение отправлено."
	UserDataFile   = "userdata.json"
	ScheduledFile  = "scheduled.json"

	// Default message
	Message_Welcome = "Добро пожаловать! Используйте команду /time для установки времени отправки сообщений."
	Message_usedStop = "Вы успешно отписались от сервиса. Ваши данные удалены."

	Wrong_timeCommand = "Использование: /time <час1> <минута1> [<час2> <минута2>]"

	// Commands
	Command_Start = "start"
	Command_Stop  = "stop"
	Command_Time  = "time"
)
