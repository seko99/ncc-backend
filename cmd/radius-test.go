package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	rad "layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type SessionsMap struct {
	sync.RWMutex
	m map[string]string
}

var radiusTestCmd = &cobra.Command{
	Use:   "radius-test",
	Short: "RADIUS Tester",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}
		log := zero.NewLogger()

		log.Info("Starting RADIUS Tester")

		log.Info("Config: %+v", cfg.Radius.Test)

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		leasesRepo := psql.NewDhcpLeases(storage, nil)

		//customers, err := customersRepo.Get()
		leases, err := leasesRepo.Get()
		if err != nil {
			panic(fmt.Sprintf("Can't get leases: %v", err))
		}

		limit := len(leases)
		if cfg.Radius.Test.Limit > 0 {
			log.Info("Limiting by %d leases", cfg.Radius.Test.Limit)
			limit = cfg.Radius.Test.Limit
		}
		accepts := 0
		rejects := 0
		total := 0
		sessions := SessionsMap{}
		sessions.m = map[string]string{}

		var nasPort uint32 = 1000000

		wg := sync.WaitGroup{}
		for _, lease := range leases {
			l := lease
			wg.Add(1)
			go func() {
				packet := rad.New(rad.CodeAccessRequest, []byte(cfg.Radius.Test.Secret))
				rfc2865.UserName_SetString(packet, l.Ip)
				rfc2865.UserPassword_SetString(packet, l.Customer.Password)
				rfc2865.NASIPAddress_Set(packet, net.ParseIP(cfg.Radius.Test.Nas.Ip))
				rfc2865.NASIdentifier_SetString(packet, cfg.Radius.Test.Nas.Identifier)
				rfc2865.NASPortType_Set(packet, rfc2865.NASPortType_Value_Ethernet)
				total++
				response, err := rad.Exchange(context.Background(), packet, cfg.Radius.Test.Auth)
				if err != nil {
					log.Error("Exchange error: %v", err)
				}
				nasPort++

				if err != nil {
					log.Error("Auth error: %v", err)
					return
				}

				if response == nil {
					log.Error("Response is nil")
					return
				}

				log.Debug("Code: %v (%s)", response.Code, l.Ip)

				if response.Code == rad.CodeAccessAccept {
					accepts++
					packet := rad.New(rad.CodeAccountingRequest, []byte(cfg.Radius.Test.Secret))
					rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_Start)
					rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
					rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
					rfc2865.UserName_SetString(packet, l.Ip)
					rfc2865.NASIPAddress_Set(packet, net.ParseIP(cfg.Radius.Test.Nas.Ip))
					rfc2865.NASIdentifier_AddString(packet, cfg.Radius.Test.Nas.Identifier)
					rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
					rfc2865.NASPort_Add(packet, rfc2865.NASPort(nasPort))
					//rfc2865.CallingStationID_Add(packet, net.ParseIP("192.168.88.239"))

					rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
					sessionId := uuid.NewString()
					rfc2866.AcctSessionID_Add(packet, []byte(sessionId))
					rfc2866.AcctSessionTime_Set(packet, 0)

					response, err = rad.Exchange(context.Background(), packet, cfg.Radius.Test.Acct)
					if err != nil {
						log.Error("Exchange error: %v", err)
					}

					sessions.Lock()
					sessions.m[l.Customer.Login] = sessionId
					sessions.Unlock()
				}
				if response.Code == rad.CodeAccessRequest {
					rejects++
				}
				wg.Done()
			}()

			limit--
			if limit <= 0 {
				break
			}
		}

		wg.Wait()
		log.Info("Accepts/rejects/total/leases: %d/%d/%d/%d", accepts, rejects, total, len(leases))

		//return

		startTime := time.Now().Unix()

		termChan := make(chan os.Signal, 1)
		signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

		go func() {
			for {
				sessions.RLock()
				s := sessions.m
				sessions.RUnlock()
				for login, sessionId := range s {
					packet := rad.New(rad.CodeAccountingRequest, []byte(cfg.Radius.Test.Secret))
					rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_InterimUpdate)
					rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
					rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
					rfc2865.UserName_SetString(packet, login)
					rfc2865.NASIPAddress_Set(packet, net.ParseIP(cfg.Radius.Test.Nas.Ip))
					rfc2865.NASIdentifier_AddString(packet, cfg.Radius.Test.Nas.Identifier)
					rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
					rfc2865.NASPort_Add(packet, rfc2865.NASPort(nasPort))
					//rfc2865.CallingStationID_Add(packet, net.ParseIP("192.168.88.239"))

					rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
					rfc2866.AcctSessionID_Add(packet, []byte(sessionId))
					rfc2866.AcctSessionTime_Set(packet, rfc2866.AcctSessionTime(time.Now().Unix()-startTime))

					response, err := rad.Exchange(context.Background(), packet, cfg.Radius.Test.Acct)
					if err != nil {
						log.Error("Exchange error: %v", err)
					}

					log.Debug("Interim response code: %s", response.Code.String())
				}
				log.Info("Sent %d interims", len(s))

				time.Sleep(cfg.Radius.Test.Interim)
			}
		}()

		<-termChan
		log.Info("Stopping Tester")

		sessions.RLock()
		s := sessions.m
		sessions.RUnlock()
		for login, sessionId := range s {
			packet := rad.New(rad.CodeAccountingRequest, []byte(cfg.Radius.Test.Secret))
			rfc2866.AcctStatusType_Set(packet, rfc2866.AcctStatusType_Value_Stop)
			rfc2865.ServiceType_Add(packet, rfc2865.ServiceType_Value_FramedUser)
			rfc2865.FramedProtocol_Add(packet, rfc2865.FramedProtocol_Value_PPP)
			rfc2865.UserName_SetString(packet, login)
			rfc2865.NASIPAddress_Add(packet, net.ParseIP(cfg.Radius.Test.Nas.Ip))
			rfc2865.NASIdentifier_AddString(packet, cfg.Radius.Test.Nas.Identifier)
			rfc2865.NASPortType_Add(packet, rfc2865.NASPortType_Value_Ethernet)
			rfc2865.NASPort_Add(packet, rfc2865.NASPort(nasPort))
			//rfc2865.CallingStationID_Add(packet, net.ParseIP("192.168.88.239"))

			rfc2866.AcctAuthentic_Add(packet, rfc2866.AcctAuthentic_Value_RADIUS)
			rfc2866.AcctSessionID_Add(packet, []byte(sessionId))
			rfc2866.AcctSessionTime_Set(packet, rfc2866.AcctSessionTime(time.Now().Unix()-startTime))
			rfc2866.AcctTerminateCause_Set(packet, rfc2866.AcctTerminateCause_Value_UserRequest)

			response, err := rad.Exchange(context.Background(), packet, cfg.Radius.Test.Acct)
			if err != nil {
				log.Error("Exchange error: %v", err)
			}

			log.Debug("Stop response code: %s", response.Code.String())
		}
	},
}
