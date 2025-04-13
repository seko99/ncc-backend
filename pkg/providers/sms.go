package providers

import "time"

//go:generate mockgen -destination=mocks/mock_sms_provider.go -package=mocks code.evixo.ru/ncc/ncc-backend/pkg/providers SmsProvider

type SmsProvider interface {
	SendOne(date time.Time, phone, message string) error
	Send(date time.Time, phones []string, message string) error
}
