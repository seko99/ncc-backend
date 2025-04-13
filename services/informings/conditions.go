package informings

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"strconv"
)

func (s *Informings) CheckConditions(customer models2.CustomerData, conds []models2.InformingConditionData) bool {
	for _, cond := range conds {
		switch cond.Field {
		case models2.FieldDeposit:
			checkResult, err := s.CheckExpression(cond, customer.Deposit)
			if err != nil || !checkResult {
				return false
			}
		case models2.FieldCredit:
			checkResult, err := s.CheckExpression(cond, customer.Credit)
			if err != nil || !checkResult {
				return false
			}
		case models2.FieldGroup:
			if customer.Group.Name != cond.Val {
				return false
			}
		case models2.FieldVerified:
			val, err := strconv.ParseBool(cond.Val)
			if err != nil {
				continue
			}

			if val && (!customer.VerifiedTs.Valid || customer.VerifiedTs.Time.IsZero()) {
				return false
			}
		case models2.FieldBlockingState:
			val, err := strconv.ParseBool(cond.Val)
			if err != nil {
				continue
			}

			if (!val && (customer.BlockingState != models2.CustomerStateActive)) ||
				(val && (customer.BlockingState == models2.CustomerStateActive)) {
				return false
			}
		case models2.FieldInternetState:
			val, err := strconv.ParseBool(cond.Val)
			if err != nil {
				continue
			}

			if (!val && (customer.ServiceInternetState != models2.ServiceStateEnabled)) ||
				(val && (customer.ServiceInternetState == models2.ServiceStateEnabled)) {
				return false
			}
		}

		for _, f := range customer.Flags {
			if f.Name == cond.Field {
				if f.Val != cond.Val {
					return false
				}
			}
		}
	}
	return true
}
