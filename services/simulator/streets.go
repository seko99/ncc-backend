package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/models"
)

func (ths *Simulator) createStreets(addresses []OsmAddress) error {
	streets := map[string]struct{}{}

	for _, a := range addresses {
		streets[a.Street] = struct{}{}
	}

	ths.log.Info("Creating %d streets...", len(streets))

	for street, _ := range streets {
		err := ths.streetsRepo.Create(models.StreetData{
			Name: street,
		})
		if err != nil {
			ths.log.Error("Can't create street: %v", err)
		}
	}

	return nil
}
