package informings

import (
	"code.evixo.ru/ncc/ncc-backend/pkg/domain"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/providers"
	repository2 "code.evixo.ru/ncc/ncc-backend/pkg/repository"
	"fmt"
	"time"
)

type Informings struct {
	log                         logger.Logger
	smsProvider                 providers.SmsProvider
	informingsRepo              repository2.Informings
	informingLogRepo            repository2.InformingLog
	informingsTestCustomersRepo repository2.InformingsTestCustomers
	customersRepo               repository2.Customers
}

type Message struct {
	Informing models2.InformingData
	Customer  models2.CustomerData
	Message   string
	Phone     string
}

func (s *Informings) Run(dryRun bool) error {

	s.log.Info("Running Informings")

	messageList, err := s.PrepareMessageList()
	if err != nil {
		return fmt.Errorf("can't prepare message list: %w", err)
	}

	if !dryRun {
		err := s.SendMessages(messageList)
		if err != nil {
			return fmt.Errorf("can't send messages: %w", err)
		}
	}

	return nil
}

func (s *Informings) ScheduledRun() error {
	return s.Run(false)
}

func (s *Informings) PrepareMessageList() ([]Message, error) {
	var messageList []Message

	informings, err := s.informingsRepo.GetEnabled()
	if err != nil {
		return nil, fmt.Errorf("can't get informings: %w", err)
	}

	customers, err := s.customersRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get customers: %w", err)
	}

	testCustomers, err := s.informingsTestCustomersRepo.Get()
	if err != nil {
		return nil, fmt.Errorf("can't get test customers: %w", err)
	}

	for _, i := range informings {

		if !i.Start.Before(time.Now()) {
			continue
		}

		customerList := customers

		if i.Mode == models2.InformingModeTest {
			customerList = []models2.CustomerData{}

			for _, c := range testCustomers {
				customerList = append(customerList, c.Customer)
			}
		}

		for _, c := range customerList {

			if !s.CheckConditions(c, i.Conditions) {
				continue
			}

			message, err := s.Replacer(i.Message, map[string]interface{}{
				"login":   c.Login,
				"deposit": c.Deposit,
				"credit":  c.Credit,
			})
			if err != nil {
				s.log.Error("can't parse message: %v", err)
				continue
			}

			s.log.Info("Prepared [%s/%s]: %s", c.Login, c.Phone, message)

			messageList = append(messageList, Message{
				Informing: models2.InformingData{
					CommonData: models2.CommonData{
						Id: i.Id,
					},
				},
				Customer: models2.CustomerData{
					CommonData: models2.CommonData{
						Id: c.Id,
					},
					Login: c.Login,
				},
				Message: message,
				Phone:   c.Phone,
			})
		}

		switch i.Repeating {
		case models2.InformingRepeatingNever:
			err := s.informingsRepo.SetState(i, models2.InformingStateDisabled)
			if err != nil {
				s.log.Error("Can't set informing state: %+v", err)
			}

			err = s.informingsRepo.SetStart(i, time.Time{})
			if err != nil {
				s.log.Error("Can't set informing start: %+v", err)
			}
		case models2.InformingRepeatingDaily:
			err := s.informingsRepo.SetStart(i, i.Start.Add(24*time.Hour))
			if err != nil {
				s.log.Error("Can't set informing start: %+v", err)
			}
		case models2.InformingRepeatingMonthly:
			err := s.informingsRepo.SetStart(i, i.Start.Add(24*time.Hour*30))
			if err != nil {
				s.log.Error("Can't set informing start: %+v", err)
			}
		}
	}

	return messageList, nil
}

func (s *Informings) SendMessages(messageList []Message) error {
	for _, m := range messageList {
		err := s.sendMessage(m.Phone, m.Message)
		if err != nil {
			s.logMessage(m, domain.MessageStatusProviderError)
			s.log.Error("Can't send message: %v", err)
		} else {
			err := s.customersRepo.SetFlag(m.Customer, models2.CustomerFlagData{
				CustomerID: models2.NewNullUUID(m.Customer.Id),
				Name:       models2.FieldSent,
				Val:        models2.ExprTrue,
			})
			if err != nil {
				s.log.Error("Can't set flag: %v", err)
			}

			s.logMessage(m, domain.MessageStatusSent)

			s.log.Info("Sent [%s/%s]: %s", m.Customer.Login, m.Phone, m.Message)
		}
	}

	return nil
}

func (s *Informings) logMessage(m Message, status int) {
	err := s.informingLogRepo.Create([]models2.InformingLogData{
		{
			CustomerId:  m.Customer.Id,
			InformingId: m.Informing.Id,
			Phone:       m.Phone,
			Message:     m.Message,
			Status:      status,
		},
	})
	if err != nil {
		s.log.Error("Can't write informing log: %+v", err)
	}
}

func (s *Informings) sendMessage(phone, message string) error {
	return s.smsProvider.SendOne(time.Now(), phone, message)
}

func NewInformings(
	log logger.Logger,
	smsProvider providers.SmsProvider,
	informingsRepo repository2.Informings,
	informingLogRepo repository2.InformingLog,
	informingsTestCustomersRepo repository2.InformingsTestCustomers,
	customersRepo repository2.Customers,
) *Informings {
	return &Informings{
		log:                         log,
		smsProvider:                 smsProvider,
		informingsRepo:              informingsRepo,
		informingLogRepo:            informingLogRepo,
		informingsTestCustomersRepo: informingsTestCustomersRepo,
		customersRepo:               customersRepo,
	}
}
