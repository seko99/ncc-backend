package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Telegram struct {
	apiKey string
	api    *tgbotapi.BotAPI
}

func (s *Telegram) Connect() error {
	api, err := tgbotapi.NewBotAPI(s.apiKey)
	if err != nil {
		return err
	}
	s.api = api
	return nil
}

func (s *Telegram) Send(chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)

	_, err := s.api.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func NewTelegram(apiKey string) *Telegram {
	return &Telegram{
		apiKey: apiKey,
	}
}
