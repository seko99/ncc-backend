package interfaces

//go:generate mockgen -destination=mocks/mock_telegram_service.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/interfaces Telegram

type Telegram interface {
	Connect() error
	Send(chatID int64, message string) error
}
