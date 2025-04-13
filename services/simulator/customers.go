package simulator

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	fixtures2 "code.evixo.ru/ncc/ncc-backend/pkg/fixtures"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"encoding/json"
	"fmt"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/fake"
	"github.com/labstack/gommon/random"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (ths *Simulator) createContracts() error {
	customers, err := ths.customerRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get customers: %w", err)
	}

	ths.log.Info("Creating %d contracts...", len(customers))

	for _, c := range customers {
		contractId := uuid.NewString()
		address := strings.Join([]string{c.City.Name, c.Street.Name, c.Build, c.Flat}, ", ")
		serial := random.New().String(4, "1234567890")
		number := serial + " " + random.New().String(8, "1234567890")
		inn := random.New().String(10, "1234567890")

		err := ths.contractsRepo.Create(models2.ContractData{
			CommonData: models2.CommonData{
				Id: contractId,
			},
			Type:               models2.ContractTypePersonal,
			Phone:              c.Phone,
			Email:              c.Email,
			Name:               c.Name,
			Document:           "паспорт",
			DocumentDate:       time.Now().AddDate(-rand.Intn(30), rand.Intn(12), rand.Intn(30)),
			DocumentIssuedBy:   "ГУ МВД России",
			Number:             number,
			Date:               time.Now(),
			Address:            address,
			ResidentialAddress: address,
			Contact:            c.Name,
			INN:                inn,
		})
		if err != nil {
			ths.log.Error("Can't create contract: %v", err)
			continue
		}

		err = ths.customerRepo.Update(models2.CustomerData{
			CommonData: models2.CommonData{
				Id: c.Id,
			},
			ContractId: models2.NewNullUUID(contractId),
		})
		if err != nil {
			ths.log.Error("Can't update customer contract: %v", err)
		}
	}

	return nil
}

func (ths *Simulator) distributeCustomers() error {
	customers, err := ths.customerRepo.Get()
	if err != nil {
		return err
	}

	ths.log.Info("Distributing %d customers...", len(customers))

	mapNodes, err := ths.mapNodesRepo.Get()
	if err != nil {
		return fmt.Errorf("can't get map nodes: %w", err)
	}

	mapNodeIdx := 0
	flat := 1

	for _, c := range customers {
		streetId := mapNodes[mapNodeIdx].StreetId
		build := mapNodes[mapNodeIdx].Build

		err := ths.customerRepo.Update(models2.CustomerData{
			CommonData: models2.CommonData{
				Id: c.Id,
			},
			StreetId: streetId,
			Build:    build,
			Flat:     strconv.Itoa(flat),
		})
		if err != nil {
			ths.log.Error("Can't update customer: %v", err)
			continue
		}

		mapNodeIdx++
		flat++
		if mapNodeIdx >= len(mapNodes) {
			mapNodeIdx = 0
			flat = 1
		}
	}

	return nil
}

func (ths *Simulator) clearCustomers() error {
	ths.log.Info("Clearing customers...")

	err := ths.clearBindings()
	if err != nil {
		return fmt.Errorf("can't clear bindings: %w", err)
	}

	r := ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_customer_fee_log")
	if r.Error != nil {
		return fmt.Errorf("can't delete fee log: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_customer_flag")
	if r.Error != nil {
		return fmt.Errorf("can't delete customer flags: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_customer_contacts")
	if r.Error != nil {
		return fmt.Errorf("can't delete customer contacts: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_payment")
	if r.Error != nil {
		return fmt.Errorf("can't delete payments: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_issue_action")
	if r.Error != nil {
		return fmt.Errorf("can't delete issue actions: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_issue")
	if r.Error != nil {
		return fmt.Errorf("can't delete issues: %w", r.Error)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_session_log")
	if r.Error != nil {
		return fmt.Errorf("can't delete session log: %w", r.Error)
	}

	err = ths.sessionsRepo.DeleteAll()
	if err != nil {
		return fmt.Errorf("can't delete sessions: %w", err)
	}

	err = ths.customerRepo.DeleteAll()
	if err != nil {
		return fmt.Errorf("can't delete customers: %w", err)
	}

	r = ths.storage.GetDB().Exec("DELETE FROM ncc.public.ncc_contract")
	if r.Error != nil {
		return fmt.Errorf("can't delete contracts: %w", r.Error)
	}

	return nil
}

func (ths *Simulator) publishEvent(eventType string, data interface{}) error {
	if ths.events != nil {
		var p map[string]interface{}

		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		err = json.Unmarshal(b, &p)
		if err != nil {
			return err
		}

		err = ths.events.PublishEvent(events.Event{
			Type:    eventType,
			Payload: p,
		})
		if err != nil {
			return fmt.Errorf("can't publish event %s: %w", eventType, err)
		}
	}

	return nil
}

func (ths *Simulator) createCustomers(num int) error {
	ths.log.Info("Creating %d customers...", num)

	usedLogins := map[string]struct{}{}

	for i := 0; i < num; i++ {

		//todo: создать договор
		//contractId := uuid.NewString()
		internet := fixtures2.FakeServicesInternet()
		serviceInternetId := internet[rand.Intn(len(internet)-1)].Id

		login := fake.UserName()

		for {
			if _, ok := usedLogins[login]; ok {
				login = fake.UserName()
				continue
			}
			usedLogins[login] = struct{}{}
			break
		}

		uid := random.New().String(10, "1234567890")
		pin := random.New().String(6, "1234567890")
		name := fmt.Sprintf("%s %s", fake.FirstName(), fake.LastName())
		deposit := rand.Intn(500)

		internetState := models2.ServiceStateEnabled

		v := rand.Intn(10)
		if v < 2 {
			deposit = rand.Intn(20) * -1
			internetState = models2.ServiceStateDisabled
		}

		err := ths.customerRepo.Create(models2.CustomerData{
			//ContractId:        models.NewNullString(contractId),
			Uid:                  uid,
			Login:                login,
			Password:             faker.Password(),
			Name:                 name,
			Phone:                faker.Phonenumber(),
			Email:                fake.EmailAddress(),
			BlockingState:        models2.CustomerStateActive,
			ServiceInternetId:    models2.NewNullUUID(serviceInternetId),
			ServiceInternetState: internetState,
			GroupId:              models2.NewNullUUID(fixtures2.FakeCustomerGroups()[0].Id),
			CityId:               models2.NewNullUUID(fixtures2.FakeCities()[0].Id),
			StreetId:             models2.NewNullUUID("00000000-0000-0000-0000-000000000000"),
			Pin:                  pin,
			Deposit:              float64(deposit),
			Credit:               0.0,
		})
		if err != nil {
			ths.log.Error("Can't create customer %s: %v", login, err)
		}

		err = ths.events.PublishEvent(events.Event{
			Type: EventTypeSimulator,
			Payload: map[string]interface{}{
				"job_type":  JobTypeCreateCustomers,
				"job_state": JobStateInProgress,
				"max":       num,
				"progress":  i,
			},
		})
		if err != nil {
			ths.log.Error("Can't publish event: %v", err)
		}
	}

	err := ths.events.PublishEvent(events.Event{
		Type: EventTypeSimulator,
		Payload: map[string]interface{}{
			"job_type":  JobTypeCreateCustomers,
			"job_state": JobStateDone,
			"max":       num,
			"progress":  num,
		},
	})
	if err != nil {
		ths.log.Error("Can't publish event: %v", err)
	}

	return nil
}
