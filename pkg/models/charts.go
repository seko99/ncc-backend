package models

import "github.com/google/uuid"

type ChartData struct {
	CommonData
	Name        string `gorm:"column:name"`
	Descr       string `gorm:"column:descr"`
	OnDashboard bool   `gorm:"column:on_dashboard"`
}

func (ChartData) TableName() string {
	return "ncc_chart"
}

type ChartParamData struct {
	CommonData
	MetricId uuid.NullUUID `gorm:"type:uuid;column:metric;not null"`
	Metric   MetricData    `gorm:"foreignKey:MetricId"`
	Color    string        `gorm:"column:color"`
	Legend   string        `gorm:"column:legend"`
	Units    string        `gorm:"column:units"`
}

func (ChartParamData) TableName() string {
	return "ncc_chart_param"
}

type ChartParamsLink struct {
	ChartId uuid.NullUUID `gorm:"type:uuid;column:chart_id;not null;index:idx_chart_param,unique"`
	ParamId uuid.NullUUID `gorm:"type:uuid;column:param_id;not null;index:idx_chart_param,unique"`
}

func (ChartParamsLink) TableName() string {
	return "ncc_chart_params_link"
}
