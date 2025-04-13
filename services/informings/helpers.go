package informings

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"strings"
)

func (s *Informings) asTemplate(field string) string {
	return fmt.Sprintf("{%s}", field)
}

func (s *Informings) Replacer(message string, data map[string]interface{}) (string, error) {
	//todo: safe templates
	msg := strings.ReplaceAll(message, s.asTemplate(models2.FieldLogin), data["login"].(string))
	msg = strings.ReplaceAll(msg, s.asTemplate(models2.FieldDeposit), fmt.Sprintf("%0.2f", data["deposit"].(float64)))
	return msg, nil
}
