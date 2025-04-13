package usecases

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/dto"
	"code.evixo.ru/ncc/ncc-backend/services/simulator/interfaces"
	"fmt"
)

type BrasGetStatUsecase struct {
	log       logger.Logger
	simulator interfaces.Simulator
}

func (ths *BrasGetStatUsecase) Execute(req dto.BrasGetStatUsecaseRequest) (dto.BrasGetStatUsecaseResponse, error) {
	sessionCache, err := ths.simulator.GetSessionCache()
	if err != nil {
		return dto.BrasGetStatUsecaseResponse{}, fmt.Errorf("can't get session: %w", err)
	}

	leasesCache, err := ths.simulator.GetLeasesCache()
	if err != nil {
		return dto.BrasGetStatUsecaseResponse{}, fmt.Errorf("can't get leases: %w", err)
	}

	customersCache, err := ths.simulator.GetCustomerCache()
	if err != nil {
		return dto.BrasGetStatUsecaseResponse{}, fmt.Errorf("can't get customers: %w", err)
	}

	nasesCache, err := ths.simulator.GetNASCache()
	if err != nil {
		return dto.BrasGetStatUsecaseResponse{}, fmt.Errorf("can't get nases: %w", err)
	}

	return dto.BrasGetStatUsecaseResponse{
		Sessions:  len(sessionCache),
		Leases:    len(leasesCache),
		Customers: len(customersCache),
		Nases:     len(nasesCache),
	}, nil
}

func NewBrasGetStatUsecase(log logger.Logger, simulator interfaces.Simulator) BrasGetStatUsecase {
	return BrasGetStatUsecase{
		log:       log,
		simulator: simulator,
	}
}
