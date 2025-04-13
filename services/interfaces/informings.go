package interfaces

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/informings"
)

//go:generate mockgen -destination=mocks/mock_informings_service.go -package=mocks code.evixo.ru/ncc/ncc-backend/services/interfaces Informings

type Informings interface {
	Run(dryRun bool) error
	CheckConditions(customer models2.CustomerData, conds []models2.InformingConditionData) bool
	Replacer(message string, data map[string]interface{}) (string, error)
	PrepareMessageList() ([]informings.Message, error)
	SendMessages(messageList []informings.Message) error
}
