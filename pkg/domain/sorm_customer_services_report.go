package domain

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"crypto/sha1"
	"fmt"
	"strings"
)

type SormCustomerServicesRecord struct {
	OrgUnit        string
	Login          string
	ContractNumber string
	ServiceCode    string
	Start          string
	End            string
	CustomData     string
}

func (s SormCustomerServicesRecord) FileName() string {
	return "abonents/services/services"
}

func (s SormCustomerServicesRecord) Header() []string {
	return []string{
		"org_unit",
		"login",
		"contract_number",
		"service_code",
		"start",
		"end",
		"custom_data",
	}
}

func (s SormCustomerServicesRecord) FromSlice(data []string) (exporter.ExportData, error) {
	var record SormCustomerServicesRecord

	if len(data) != 7 {
		return record, fmt.Errorf("wrong data: %v", data)
	}

	record.OrgUnit = data[0]
	record.Login = data[1]
	record.ContractNumber = data[2]
	record.ServiceCode = data[3]
	record.Start = data[4]
	record.End = data[5]
	record.CustomData = data[6]
	return record, nil
}

func (s SormCustomerServicesRecord) ToSlice() []string {
	return []string{
		s.OrgUnit,
		s.Login,
		s.ContractNumber,
		s.ServiceCode,
		s.Start,
		s.End,
		s.CustomData,
	}
}

func (s SormCustomerServicesRecord) GetHash() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join([]string{
		s.OrgUnit,
		s.Login,
		s.ContractNumber,
		s.ServiceCode,
		s.Start,
		s.End,
		s.CustomData,
	}, ","))))
}

func (s SormCustomerServicesRecord) ToSormCustomerServiceData() models.SormCustomerServiceData {
	return models.SormCustomerServiceData{
		Hash:           s.GetHash(),
		OrgUnit:        s.OrgUnit,
		Login:          s.Login,
		ContractNumber: s.ContractNumber,
		ServiceCode:    s.ServiceCode,
		Start:          s.Start,
		End:            s.End,
		CustomData:     s.CustomData,
	}
}
