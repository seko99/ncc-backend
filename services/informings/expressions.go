package informings

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"strconv"
)

func (s *Informings) CheckExpression(cond models2.InformingConditionData, value interface{}) (bool, error) {
	switch v := value.(type) {
	case float64:
		val, err := strconv.ParseFloat(cond.Val, 64)
		if err != nil {
			return false, err
		}

		switch cond.Expr {
		case models2.ExprEq:
			if v != val {
				return false, nil
			}
		case models2.ExprGt:
			if v <= val {
				return false, nil
			}
		case models2.ExprGe:
			if v < val {
				return false, nil
			}
		case models2.ExprLt:
			if v >= val {
				return false, nil
			}
		case models2.ExprLe:
			if v > val {
				return false, nil
			}
		}
	case string:
		switch cond.Expr {
		case models2.ExprEq:
			if v != cond.Val {
				return false, nil
			}
		case models2.ExprNe:
			if v == cond.Val {
				return false, nil
			}
		}
	}

	return true, nil
}
