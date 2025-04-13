package dto

import "code.evixo.ru/ncc/ncc-backend/pkg/models"

type BrasGetSessionsUsecaseRequest struct{}

type BrasGetSessionsUsecaseResponse struct {
	Sessions []models.SessionData
}
