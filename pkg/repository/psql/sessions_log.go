package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"fmt"
	"gorm.io/gorm/clause"
)

type SessionsLog struct {
	storage *psqlstorage.Storage
}

func (s *SessionsLog) Create(session models.SessionsLogData) error {
	r := s.storage.GetDB().Model(&models.SessionsLogData{}).
		Create(&session)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SessionsLog) DeleteById(id string) error {
	r := s.storage.GetDB().Model(&models.SessionsLogData{}).
		Delete("id = @id", sql.Named("id", id))
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *SessionsLog) GetBySessionId(sessionId string) (models.SessionsLogData, error) {
	var session models.SessionsLogData

	r := s.storage.GetDB().Model(&models.SessionsLogData{}).
		Preload(clause.Associations).
		Where("acct_session_id = @session_id", sql.Named("session_id", sessionId)).
		First(&session)
	if r.Error != nil {
		return models.SessionsLogData{}, r.Error
	}

	return session, nil
}

func (s *SessionsLog) GetByCustomer(id string, period repository.TimePeriod, limit ...int) ([]models.SessionsLogData, error) {
	var sessions []models.SessionsLogData

	periodClause := repository.PeriodClause("start_time", period)
	if !period.In.IsZero() {
		start := fmt.Sprintf("%d-%d-%d 00:00:00", period.In.Year(), period.In.Month(), period.In.Day())
		end := fmt.Sprintf("%d-%d-%d 23:59:59", period.In.Year(), period.In.Month(), period.In.Day())
		periodClause = fmt.Sprintf(`(start_time>='%s' AND start_time<='%s') 
			OR (stop_time>='%s' AND stop_time<='%s')
			OR (start_time<'%s' AND stop_time>'%s')`, start, end, start, end, start, end)
	}

	r := s.storage.GetDB().Model(&models.SessionsLogData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("customer_id = @id", sql.Named("id", id)).
		Where(periodClause).
		Order("start_time")

	if len(limit) > 0 {
		r = r.Limit(limit[0])
	}

	r = r.Find(&sessions)

	if r.Error != nil {
		return nil, r.Error
	}

	return sessions, nil
}

func NewSessionsLog(storage *psqlstorage.Storage) *SessionsLog {
	return &SessionsLog{
		storage: storage,
	}
}
