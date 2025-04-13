package cmd

import (
	"bufio"
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger/zero"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"code.evixo.ru/ncc/ncc-backend/pkg/storage/psql"
	"fmt"
	"github.com/spf13/cobra"
	dblogger "gorm.io/gorm/logger"
	"os"
	"strings"
)

func askForConfirmation(s string, yes ...bool) bool {
	if len(yes) == 1 && yes[0] {
		return true
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

var (
	modelsList = map[string]interface{}{
		"snapshots": &models2.SnapshotData{},
	}
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "DB migration",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.NewConfig()
		if err != nil {
			panic(err)
		}
		log := zero.NewLogger()

		log.Info("Starting migration")

		storage := psqlstorage.NewStorage(cfg, log, psqlstorage.WithLogLevel(dblogger.Error))
		err = storage.Connect()
		if err != nil {
			panic(fmt.Sprintf("can't connect to storage: %v", err))
		}

		modelFlag, _ := cmd.Flags().GetString("model")
		yesFlag, _ := cmd.Flags().GetBool("yes")

		if modelFlag != "" {
			if model, ok := modelsList[modelFlag]; ok {
				if askForConfirmation(fmt.Sprintf("Migrate model %s?", modelFlag), yesFlag) {
					err = storage.GetDB().AutoMigrate(model)
				}
			}

			return
		}

		if askForConfirmation("Migrate all models?", yesFlag) {
			err := storage.Migrate()
			if err != nil {
				log.Error("Migration error: %v", err)
				os.Exit(1)
			}
		}
	},
}
