package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"errors"
	"gorm.io/gorm"
)

type Preferences struct {
	storage *psqlstorage.Storage
}

func NewPreferences(storage *psqlstorage.Storage) *Preferences {
	return &Preferences{
		storage: storage,
	}
}

func (p *Preferences) GetFeeProcessingInProgress() (bool, error) {
	var preference models.PreferencesData

	err := p.storage.GetDB().
		Where("name = ?", "feeProcessingInProgress").
		First(&preference).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return preference.BooleanValue, nil
}

func (p *Preferences) SetFeeProcessingInProgress(value bool) error {
	var preference models.PreferencesData

	err := p.storage.GetDB().
		Where("name = ?", "feeProcessingInProgress").
		First(&preference).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			preference = models.PreferencesData{
				Name:         "feeProcessingInProgress",
				BooleanValue: value,
			}
			return p.storage.GetDB().Create(&preference).Error
		}
		return err
	}

	preference.BooleanValue = value
	return p.storage.GetDB().Save(&preference).Error
}
