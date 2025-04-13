package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"time"
)

func (ths *Simulator) DeleteIssues() error {
	err := ths.issuesRepo.DeleteAll()
	if err != nil {
		return err
	}
	return nil
}

func (ths *Simulator) CreateIssues() error {
	ths.log.Info("Creating issues...")

	customers, err := ths.customerRepo.Get(20)
	if err != nil {
		return fmt.Errorf("can't get customers: %w", err)
	}

	iid := 100_000
	for _, c := range customers {
		err := ths.issuesRepo.Create(models.IssueData{
			IId:         iid,
			Date:        time.Now(),
			Status:      models.IssueStatusOpen,
			CustomerId:  models.NewNullUUID(c.Id),
			IssueTypeId: models.NewNullUUID(fixtures.IssueTypeFailure),
			UrgencyId:   models.NewNullUUID(fixtures.IssueUrgencyMedium),
			CityId:      c.CityId,
			StreetId:    c.StreetId,
			Build:       c.Build,
			Flat:        c.Flat,
			Phone:       c.Phone,
			Name:        c.Name,
		})
		if err != nil {
			ths.log.Error("Can't create issue: %v", err)
			continue
		}

		iid++
	}

	return nil
}
