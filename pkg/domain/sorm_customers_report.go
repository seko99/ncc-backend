package domain

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/services/interfaces/exporter"
	"crypto/sha1"
	"fmt"
	"strings"
)

type SormCustomersRecord struct {
	OrgUnit               string
	Login                 string
	IP                    string
	Email                 string
	Phone                 string
	MAC                   string
	ContractDate          string
	ContractNumber        string
	Status                string
	Start                 string
	End                   string
	Type                  string
	NameType              string
	FirstName             string
	MiddleName            string
	LastName              string
	Name                  string
	BirthDate             string
	DocType               string
	DocSerial             string
	DocNumber             string
	DocIssuer             string
	Doc                   string
	DocCode               string
	Bank                  string
	BankAccount           string
	EnterpriseName        string
	EnterpriseINN         string
	ContactName           string
	ContactPhone          string
	EnterpriseBank        string
	EnterpriseBankAccount string
	RegistrationType      string
	RegistrationZip       string
	RegistrationCountry   string
	RegistrationRegion    string
	RegistrationDistrict  string
	RegistrationCity      string
	RegistrationStreet    string
	RegistrationBuild     string
	RegistrationCorp      string
	RegistrationFlat      string
	RegistrationAddress   string
	DeviceType            string
	DeviceZip             string
	DeviceCountry         string
	DeviceRegion          string
	DeviceDistrict        string
	DeviceCity            string
	DeviceStreet          string
	DeviceBuild           string
	DeviceCorp            string
	DeviceFlat            string
	DeviceAddress         string
	PostType              string
	PostZip               string
	PostCountry           string
	PostRegion            string
	PostDistrict          string
	PostCity              string
	PostStreet            string
	PostBuild             string
	PostCorp              string
	PostFlat              string
	PostAddress           string
	BillType              string
	BillZip               string
	BillCountry           string
	BillRegion            string
	BillDistrict          string
	BillCity              string
	BillStreet            string
	BillBuild             string
	BillCorp              string
	BillFlat              string
	BillAddress           string
}

func (s SormCustomersRecord) FileName() string {
	return "abonents/abonents/abonents_new"
}

func (s SormCustomersRecord) Header() []string {
	return []string{}
}

func (s SormCustomersRecord) FromSlice(data []string) (exporter.ExportData, error) {
	var record SormCustomersRecord

	if len(data) != 76 {
		return record, fmt.Errorf("wrong data: %v", data)
	}

	record.OrgUnit = data[0]
	record.Login = data[1]
	record.IP = data[2]
	record.Email = data[3]
	record.Phone = data[4]
	record.MAC = data[5]
	record.ContractDate = data[6]
	record.ContractNumber = data[7]
	record.Status = data[8]
	record.Start = data[9]
	record.End = data[10]
	record.Type = data[11]
	record.NameType = data[12]
	record.FirstName = data[13]
	record.MiddleName = data[14]
	record.LastName = data[15]
	record.Name = data[16]
	record.BirthDate = data[17]
	record.DocType = data[18]
	record.DocSerial = data[19]
	record.DocNumber = data[20]
	record.DocIssuer = data[21]
	record.Doc = data[22]
	record.DocCode = data[23]
	record.Bank = data[24]
	record.BankAccount = data[25]
	record.EnterpriseName = data[26]
	record.EnterpriseINN = data[27]
	record.ContactName = data[28]
	record.ContactPhone = data[29]
	record.EnterpriseBank = data[30]
	record.EnterpriseBankAccount = data[31]
	record.RegistrationType = data[32]
	record.RegistrationZip = data[33]
	record.RegistrationCountry = data[34]
	record.RegistrationRegion = data[35]
	record.RegistrationDistrict = data[36]
	record.RegistrationCity = data[37]
	record.RegistrationStreet = data[38]
	record.RegistrationBuild = data[39]
	record.RegistrationCorp = data[40]
	record.RegistrationFlat = data[41]
	record.RegistrationAddress = data[42]
	record.DeviceType = data[43]
	record.DeviceZip = data[44]
	record.DeviceCountry = data[45]
	record.DeviceRegion = data[46]
	record.DeviceDistrict = data[47]
	record.DeviceCity = data[48]
	record.DeviceStreet = data[49]
	record.DeviceBuild = data[50]
	record.DeviceCorp = data[51]
	record.DeviceFlat = data[52]
	record.DeviceAddress = data[53]
	record.PostType = data[54]
	record.PostZip = data[55]
	record.PostCountry = data[56]
	record.PostRegion = data[57]
	record.PostDistrict = data[58]
	record.PostCity = data[59]
	record.PostStreet = data[60]
	record.PostBuild = data[61]
	record.PostCorp = data[62]
	record.PostFlat = data[63]
	record.PostAddress = data[64]
	record.BillType = data[65]
	record.BillZip = data[66]
	record.BillCountry = data[67]
	record.BillRegion = data[68]
	record.BillDistrict = data[69]
	record.BillCity = data[70]
	record.BillStreet = data[71]
	record.BillBuild = data[72]
	record.BillCorp = data[73]
	record.BillFlat = data[74]
	record.BillAddress = data[75]
	return record, nil
}

func (s SormCustomersRecord) ToSlice() []string {
	return []string{
		s.OrgUnit,
		s.Login,
		s.IP,
		s.Email,
		s.Phone,
		s.MAC,
		s.ContractDate,
		s.ContractNumber,
		s.Status,
		s.Start,
		s.End,
		s.Type,
		s.NameType,
		s.FirstName,
		s.MiddleName,
		s.LastName,
		s.Name,
		s.BirthDate,
		s.DocType,
		s.DocSerial,
		s.DocNumber,
		s.DocIssuer,
		s.Doc,
		s.DocCode,
		s.Bank,
		s.BankAccount,
		s.EnterpriseName,
		s.EnterpriseINN,
		s.ContactName,
		s.ContactPhone,
		s.EnterpriseBank,
		s.EnterpriseBankAccount,
		s.RegistrationType,
		s.RegistrationZip,
		s.RegistrationCountry,
		s.RegistrationRegion,
		s.RegistrationDistrict,
		s.RegistrationCity,
		s.RegistrationStreet,
		s.RegistrationBuild,
		s.RegistrationCorp,
		s.RegistrationFlat,
		s.RegistrationAddress,
		s.DeviceType,
		s.DeviceZip,
		s.DeviceCountry,
		s.DeviceRegion,
		s.DeviceDistrict,
		s.DeviceCity,
		s.DeviceStreet,
		s.DeviceBuild,
		s.DeviceCorp,
		s.DeviceFlat,
		s.DeviceAddress,
		s.PostType,
		s.PostZip,
		s.PostCountry,
		s.PostRegion,
		s.PostDistrict,
		s.PostCity,
		s.PostStreet,
		s.PostBuild,
		s.PostCorp,
		s.PostFlat,
		s.PostAddress,
		s.BillType,
		s.BillZip,
		s.BillCountry,
		s.BillRegion,
		s.BillDistrict,
		s.BillCity,
		s.BillStreet,
		s.BillBuild,
		s.BillCorp,
		s.BillFlat,
		s.BillAddress,
	}
}

func (s SormCustomersRecord) GetHash() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(strings.Join([]string{
		s.OrgUnit,
		s.Login,
		s.IP,
		s.Email,
		s.Phone,
		s.MAC,
		s.ContractDate,
		s.ContractNumber,
		s.Status,
		s.Start,
		s.End,
		s.Type,
		s.NameType,
		s.FirstName,
		s.MiddleName,
		s.LastName,
		s.Name,
		s.BirthDate,
		s.DocType,
		s.DocSerial,
		s.DocNumber,
		s.DocIssuer,
		s.Doc,
		s.DocCode,
		s.Bank,
		s.BankAccount,
		s.EnterpriseName,
		s.EnterpriseINN,
		s.ContactName,
		s.ContactPhone,
		s.EnterpriseBank,
		s.EnterpriseBankAccount,
		s.RegistrationType,
		s.RegistrationZip,
		s.RegistrationCountry,
		s.RegistrationRegion,
		s.RegistrationDistrict,
		s.RegistrationCity,
		s.RegistrationStreet,
		s.RegistrationBuild,
		s.RegistrationCorp,
		s.RegistrationFlat,
		s.RegistrationAddress,
		s.DeviceType,
		s.DeviceZip,
		s.DeviceCountry,
		s.DeviceRegion,
		s.DeviceDistrict,
		s.DeviceCity,
		s.DeviceStreet,
		s.DeviceBuild,
		s.DeviceCorp,
		s.DeviceFlat,
		s.DeviceAddress,
		s.PostType,
		s.PostZip,
		s.PostCountry,
		s.PostRegion,
		s.PostDistrict,
		s.PostCity,
		s.PostStreet,
		s.PostBuild,
		s.PostCorp,
		s.PostFlat,
		s.PostAddress,
		s.BillType,
		s.BillZip,
		s.BillCountry,
		s.BillRegion,
		s.BillDistrict,
		s.BillCity,
		s.BillStreet,
		s.BillBuild,
		s.BillCorp,
		s.BillFlat,
		s.BillAddress,
	}, ","))))
}

func (s SormCustomersRecord) ToSormCustomersData() models.SormCustomersData {
	return models.SormCustomersData{
		Hash:                  s.GetHash(),
		OrgUnit:               s.OrgUnit,
		Login:                 s.Login,
		IP:                    s.IP,
		Email:                 s.Email,
		Phone:                 s.Phone,
		MAC:                   s.MAC,
		ContractDate:          s.ContractDate,
		ContractNumber:        s.ContractNumber,
		Status:                s.Status,
		Start:                 s.Start,
		End:                   s.End,
		Type:                  s.Type,
		NameType:              s.NameType,
		FirstName:             s.FirstName,
		MiddleName:            s.MiddleName,
		LastName:              s.LastName,
		Name:                  s.Name,
		BirthDate:             s.BirthDate,
		DocType:               s.DocType,
		DocSerial:             s.DocSerial,
		DocNumber:             s.DocNumber,
		DocIssuer:             s.DocIssuer,
		Doc:                   s.Doc,
		DocCode:               s.DocCode,
		Bank:                  s.Bank,
		BankAccount:           s.BankAccount,
		EnterpriseName:        s.EnterpriseName,
		EnterpriseINN:         s.EnterpriseINN,
		ContactName:           s.ContactName,
		ContactPhone:          s.ContactPhone,
		EnterpriseBank:        s.EnterpriseBank,
		EnterpriseBankAccount: s.EnterpriseBankAccount,
		RegistrationType:      s.RegistrationType,
		RegistrationZip:       s.RegistrationZip,
		RegistrationCountry:   s.RegistrationCountry,
		RegistrationRegion:    s.RegistrationRegion,
		RegistrationDistrict:  s.RegistrationDistrict,
		RegistrationCity:      s.RegistrationCity,
		RegistrationStreet:    s.RegistrationStreet,
		RegistrationBuild:     s.RegistrationBuild,
		RegistrationCorp:      s.RegistrationCorp,
		RegistrationFlat:      s.RegistrationFlat,
		RegistrationAddress:   s.RegistrationAddress,
		DeviceType:            s.DeviceType,
		DeviceZip:             s.DeviceZip,
		DeviceCountry:         s.DeviceCountry,
		DeviceRegion:          s.DeviceRegion,
		DeviceDistrict:        s.DeviceDistrict,
		DeviceCity:            s.DeviceCity,
		DeviceStreet:          s.DeviceStreet,
		DeviceBuild:           s.DeviceBuild,
		DeviceCorp:            s.DeviceCorp,
		DeviceFlat:            s.DeviceFlat,
		DeviceAddress:         s.DeviceAddress,
		PostType:              s.PostType,
		PostZip:               s.PostZip,
		PostCountry:           s.PostCountry,
		PostRegion:            s.PostRegion,
		PostDistrict:          s.PostDistrict,
		PostCity:              s.PostCity,
		PostStreet:            s.PostStreet,
		PostBuild:             s.PostBuild,
		PostCorp:              s.PostCorp,
		PostFlat:              s.PostFlat,
		PostAddress:           s.PostAddress,
		BillType:              s.BillType,
		BillZip:               s.BillZip,
		BillCountry:           s.BillCountry,
		BillRegion:            s.BillRegion,
		BillDistrict:          s.BillDistrict,
		BillCity:              s.BillCity,
		BillStreet:            s.BillStreet,
		BillBuild:             s.BillBuild,
		BillCorp:              s.BillCorp,
		BillFlat:              s.BillFlat,
		BillAddress:           s.BillAddress,
	}
}
