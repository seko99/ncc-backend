package providers

import (
	"bytes"
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SmsRequest struct {
	DateSend   string   `json:"dateSend"`
	Message    string   `json:"message"`
	PhonesList []string `json:"phonesList"`
}

type PhoenixSms struct {
	cfg config.SmsProviderConfig
	log logger.Logger
}

func (s *PhoenixSms) SendOne(date time.Time, phone, message string) error {
	return s.Send(date, []string{phone}, message)
}

func (s *PhoenixSms) Send(date time.Time, phones []string, message string) error {
	messageDate := date.Add(5 * time.Minute).Format("02.01.2006 15:04:05")
	body, err := json.Marshal(SmsRequest{
		DateSend:   messageDate,
		Message:    message,
		PhonesList: phones,
	})
	if err != nil {
		return fmt.Errorf("can't marshal sms request: %w", err)
	}

	br := bytes.NewReader(body)

	url := s.cfg.URI + "/dispatches?token=" + s.cfg.Token
	req, err := http.NewRequest(http.MethodPost, url, br)
	if err != nil {
		return fmt.Errorf("can't form request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("can't send request: %w", err)
	}

	s.log.Info("Sending message dispatch for %v to %v", messageDate, phones)

	if response.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(response.Body)
		return fmt.Errorf("PhoenixSms error: %d %s", response.StatusCode, string(b))
	}

	return nil
}

func NewPhoenixSms(
	cfg config.SmsProviderConfig,
	log logger.Logger,
) *PhoenixSms {
	return &PhoenixSms{
		cfg: cfg,
		log: log,
	}
}
