package data

import (
// user "Bot_Compliment/userdata"
)

var (
	DEFAUTLE_MESSAGE = "Ваше отложенное сообщение отправлено."
	USER_DATA_FILE   = "userdata.json"
	SCHEDULED_FILE  = "scheduled.json"

	// Default message
	MESSAGE_WELLCOME = "Добро пожаловать! Используйте команду /time для установки времени отправки сообщений."
	MESSAGE_STOP = "Вы успешно отписались от сервиса. Ваши данные удалены."

	WRONG_TIME_COMMAND = "Использование: /time <час1> <минута1> [<час2> <минута2>]"

	// Commands
	COMMAND_START = "start"
	COMMAND_STOP  = "stop"
	COMMAND_SET_TIME = "time"
)
