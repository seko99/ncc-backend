package simulator

import (
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"github.com/labstack/gommon/random"
	"time"
)

func (ths *Simulator) UpdateSessions() error {

	sessions, err := ths.sessionsRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get sessions: %w", err)
	}

	ths.log.Info("Updating %d sessions...", len(sessions))

	for _, session := range sessions {
		session.Duration = int64(time.Since(session.StartTime).Seconds())
		session.LastAlive = time.Now()

		err := ths.sessionsRepo.Update(session)
		if err != nil {
			ths.log.Error("Can't update session: %v", err)
		}
	}

	return nil
}

func (ths *Simulator) disconnectSessions() error {
	ths.log.Info("Disconnecting sessions...")

	sessions, err := ths.sessionsRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get sessions: %w", err)
	}

	for _, session := range sessions {
		if session.Customer.ServiceInternetState != models2.ServiceStateEnabled {
			err := ths.sessionsRepo.Delete(session.Id)
			if err != nil {
				ths.log.Error("Can't delete session: %v", err)
			}
		}
	}

	return nil
}

func (ths *Simulator) createSessions() error {
	leases, err := ths.dhcpLeasesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get leases: %w", err)
	}

	nases, err := ths.nasesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get nases: %w", err)
	}

	ths.log.Info("Creating %d sessions...", len(leases))

	for _, l := range leases {
		if l.Customer.ServiceInternetState != models2.ServiceStateEnabled {
			continue
		}

		err = ths.sessionsRepo.Create([]models2.SessionData{
			{
				AcctSessionId:     random.New().String(12, "0123456789abcdef"),
				Login:             l.Customer.Login,
				CustomerId:        l.CustomerId,
				StartTime:         time.Now().Add(-3 * time.Hour),
				Duration:          3 * 60 * 60,
				ServiceInternetId: l.Customer.ServiceInternetId,
				LastAlive:         time.Now(),
				Ip:                l.Ip,
				Mac:               l.Mac,
				Remote:            l.Remote,
				NasId:             models2.NewNullUUID(nases[0].Id),
				NasName:           nases[0].Name,
			},
		})
		if err != nil {
			ths.log.Error("Can't create session: %v", err)
		}
	}

	return nil
}

func (ths *Simulator) DropSessions() error {
	return ths.sessionsRepo.DeleteAll()
}
