package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/events"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	"code.evixo.ru/ncc/ncc-backend/services/dhcp"
	"code.evixo.ru/ncc/ncc-backend/services/radius"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var dhcpCmd = &cobra.Command{
	Use:   "dhcp-server",
	Short: "DHCP server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}
		log := zero.NewLogger()

		log.Info("Initializing event system...")
		events, err := events.NewEvents(cfg, log, uuid.NewString(), radius.Queue)
		if err != nil {
			panic(fmt.Sprintf("Can't init event system: %v", err))
		}

		log.Info("Event system initialized")

		log.Info("Starting DHCP server...")

		dhcpServer := dhcp.NewDhcpServer(cfg, log, events)

		err = dhcpServer.Start()
		if err != nil {
			panic(err)
		}
	},
}
