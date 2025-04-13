package repository

import (
	"fmt"
	"strings"
	"time"
)

type TimePeriod struct {
	Start       time.Time
	End         time.Time
	In          time.Time
	Excluding   bool
	ExtraFields []string
}

func PeriodClause(field string, period TimePeriod) string {
	equal := ""
	if !period.Excluding {
		equal = "="
	}

	// between
	if !period.Start.IsZero() && !period.End.IsZero() {
		return fmt.Sprintf("%s >%s '%s' AND %s <%s '%s'", field, equal, period.Start.Format("2006-01-02 15:04:05"), field, equal, period.End.Format("2006-01-02 15:04:05"))
	}

	// starts
	if !period.Start.IsZero() {
		fieldClauses := []string{fmt.Sprintf("%s >%s '%s'", field, equal, period.Start.Format("2006-01-02 15:04:05"))}
		for _, f := range period.ExtraFields {
			fieldClauses = append(fieldClauses, fmt.Sprintf("%s >%s '%s'", f, equal, period.Start.Format("2006-01-02 15:04:05")))
		}
		return strings.Join(fieldClauses, " OR ")
	}

	// ends
	if !period.End.IsZero() {
		return fmt.Sprintf("%s <%s '%s'", field, equal, period.End.Format("2006-01-02 15:04:05"))
	}

	return ""
}
