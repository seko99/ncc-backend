package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/handlers"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	psql2 "code.evixo.ru/ncc/ncc-backend/pkg/repository/psql"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"time"
)

const (
	apiPrefix = "/v1"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "API daemon",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}

		log := zero.NewLogger(zerolog.DebugLevel)

		log.Info("Starting API daemon")

		storage := psqlstorage.NewStorage(cfg, log)
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		go func() {
			for {
				time.Sleep(cfg.Watcher.Delay)

			}
		}()

		g := gin.New()

		customersRepo := psql2.NewCustomers(storage, nil)
		paymentsRepo := psql2.NewPayments(storage)
		/*		authRepo := psql.NewAuth(cfg, storage)
				feeRepo := psql.NewFees(storage)
		*/

		reportsEndpoint := handlers.NewReports(customersRepo, paymentsRepo)

		reportsGroup := g.Group("/reports")
		{
			reportsGroup.GET("/payments_by_days/:start/:end", reportsEndpoint.Execute())
		}

		err = g.Run(cfg.Listen)
		if err != nil {
			log.Error("can't listen: %v", err)
		}

	},
}
