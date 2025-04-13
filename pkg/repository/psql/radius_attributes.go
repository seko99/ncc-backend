package psql

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"database/sql"
	"gorm.io/gorm/clause"
	"time"
)

type RadiusAttributes struct {
	storage *psqlstorage.Storage
}

func (s *RadiusAttributes) Create(data models.RadiusAttributeData) error {
	r := s.storage.GetDB().Model(&models.RadiusAttributeData{}).Create(&data)
	if r.Error != nil {
		return r.Error
	}

	r = s.storage.GetDB().Model(&models.RadiusAttributeLink{}).
		Create(&models.RadiusAttributeLink{
			VendorId:    data.VendorId,
			AttributeId: models.NewNullUUID(data.Id),
		})
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *RadiusAttributes) Upsert(data models.RadiusAttributeData) error {
	r := s.storage.GetDB().Model(&models.RadiusAttributeData{}).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&data)
	if r.Error != nil {
		return r.Error
	}

	r = s.storage.GetDB().Model(&models.RadiusAttributeLink{}).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "vendor_id"}, {Name: "attribute_id"}},
			DoNothing: true,
		}).
		Create(&models.RadiusAttributeLink{
			VendorId:    data.VendorId,
			AttributeId: models.NewNullUUID(data.Id),
		})
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *RadiusAttributes) Update(data models.RadiusAttributeData) error {
	r := s.storage.GetDB().Model(&models.RadiusAttributeData{}).
		Where("id = @id", sql.Named("id", data.Id)).
		Updates(data)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

func (s *RadiusAttributes) Delete(id string) error {
	r := s.storage.GetDB().Model(&models.RadiusAttributeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Where("id = @id", sql.Named("id", id)).
		Update("delete_ts", time.Now())
	if r.Error != nil {
		return r.Error
	}

	return nil
}

func (s *RadiusAttributes) Get() ([]models.RadiusAttributeData, error) {
	var data []models.RadiusAttributeData

	r := s.storage.GetDB().Model(&models.RadiusAttributeData{}).
		Preload(clause.Associations).
		Where("delete_ts is null").
		Find(&data)
	if r.Error != nil {
		return nil, r.Error
	}

	return data, nil
}

func NewRadiusAttributes(storage *psqlstorage.Storage) *RadiusAttributes {
	return &RadiusAttributes{
		storage: storage,
	}
}
