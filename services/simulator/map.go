package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
)

func (ths *Simulator) UpdateMap() error {
	ths.log.Info("Updating map...")

	mapNodes, err := ths.mapNodesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get nap nodes: %w", err)
	}

	sessions, err := ths.sessionsRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get sessions: %w", err)
	}

	for _, m := range mapNodes {
		newStatus := models2.MapNodeStatusInactive

		for _, session := range sessions {
			if session.Customer.CityId == m.CityId &&
				session.Customer.StreetId == m.StreetId &&
				session.Customer.Build == m.Build {

				newStatus = models2.MapNodeStatusActive
				break
			}
		}

		if newStatus != m.Status {
			err := ths.mapNodesRepo.Update(models2.MapNodeData{
				CommonData: models2.CommonData{
					Id: m.Id,
				},
				Status: newStatus,
			})
			if err != nil {
				ths.log.Error("Can't update map node: %v", err)
			}
		}
	}

	return nil
}

func (ths *Simulator) clearMapNodes() error {
	ths.log.Info("Clearing map nodes...")

	r := ths.storage.GetDB().Exec("UPDATE ncc.public.ncc_device SET map_node_id=NULL")
	if r.Error != nil {
		return fmt.Errorf("can't set NULL to devices map_node: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_map_node")
	if r.Error != nil {
		return fmt.Errorf("can't delete map nodes: %w", r.Error)
	}

	return nil
}

func (ths *Simulator) createMapNodes(addresses []OsmAddress) error {
	ths.log.Info("Creating %d map nodes...", len(addresses))

	streets, err := ths.streetsRepo.Get()
	if err != nil {
		return err
	}

	cities, err := ths.citiesRepo.Get()
	if err != nil {
		return err
	}

	for i, a := range addresses {
		streetId := ""
		for _, street := range streets {
			if street.Name == a.Street {
				streetId = street.Id
				break
			}
		}
		err = ths.mapNodesRepo.Create(models2.MapNodeData{
			Lat:      a.LatLng.Lat,
			Lng:      a.LatLng.Lng,
			CityId:   models2.NewNullUUID(cities[0].Id),
			StreetId: models2.NewNullUUID(streetId),
			Build:    a.Build,
			Status:   models2.MapNodeStatusInactive,
		})
		if err != nil {
			ths.log.Error("Can't create map node for %s: %v", a.Id, err)
		}

		err = ths.events.PublishEvent(events.Event{
			Type: EventTypeSimulator,
			Payload: map[string]interface{}{
				"job_type":  JobTypeCreateMapNodes,
				"job_state": JobStateInProgress,
				"max":       len(addresses),
				"progress":  i,
			},
		})
		if err != nil {
			ths.log.Error("Can't publish event: %v", err)
		}
	}

	err = ths.events.PublishEvent(events.Event{
		Type: EventTypeSimulator,
		Payload: map[string]interface{}{
			"job_type":  JobTypeCreateMapNodes,
			"job_state": JobStateDone,
			"max":       len(addresses),
			"progress":  len(addresses),
		},
	})
	if err != nil {
		ths.log.Error("Can't publish event: %v", err)
	}

	return nil
}
