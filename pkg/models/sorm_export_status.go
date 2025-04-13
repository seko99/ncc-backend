package models

type SormExportStatusData struct {
	CommonData

	FileName    string `gorm:"column:file_name;unique"`
	ExportCount int    `gorm:"column:export_count"`
	Errors      int    `gorm:"column:errors"`
	Status      string `gorm:"column:status"`
}

func (SormExportStatusData) TableName() string {
	return "ncc_sorm_export_status"
}
