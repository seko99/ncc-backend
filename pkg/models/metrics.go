package models

const (
	MetricValueTypeDouble = "DOUBLE"
	MetricValueTypeLong   = "LONG"
)

type MetricData struct {
	CommonData
	Name      string
	ValueType string
}

func (MetricData) TableName() string {
	return "ncc_metric"
}
