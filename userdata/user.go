package userdata

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
	UserDatas []UserData
)