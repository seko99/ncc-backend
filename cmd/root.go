package cmd

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/utils"
	"github.com/spf13/cobra"
	"time"
)

const (
	prefix = "/v1"
)

var (
	cfgFile   string
	debugMode *bool
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "NetEngine",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize()

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "configs/local.yaml", "config file")
	debugMode = rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode")

	nowString := time.Now().Format("2006-01-02 15:04:05")
	rootCmd.PersistentFlags().BoolP("daemon", "", false, "--daemon")
	rootCmd.PersistentFlags().StringP("time", "t", nowString, "--time='2022-09-03 04:00:00'")
	rootCmd.PersistentFlags().BoolP("dry", "", false, "--dry")
	rootCmd.PersistentFlags().BoolP("fix", "f", false, "--fix")
	rootCmd.PersistentFlags().BoolP("yes", "y", false, "--yes")
	rootCmd.PersistentFlags().BoolP("ignore-credit-expire", "", false, "--ignore-credit-expire")
	rootCmd.PersistentFlags().BoolP("ignore-service-state", "", false, "--ignore-service-state")
	rootCmd.PersistentFlags().StringP("model", "m", "", "--model=snapshots")
	rootCmd.PersistentFlags().StringP("login", "l", "", "--login=admin")
	rootCmd.PersistentFlags().StringP("uid", "", "", "--uid=4423,3883,123")
	rootCmd.PersistentFlags().IntP("max-blocks", "b", 0, "--max-blocks=100")
	rootCmd.PersistentFlags().BoolP("zero-fee", "", false, "--zero-fee")
	rootCmd.PersistentFlags().StringP("session-stop", "", "", "--session-stop=2022-09-21")

	// migrator flags
	{
		migrateCmd.PersistentFlags().StringP("model", "m", "", "--model=snapshots")
		migrateCmd.PersistentFlags().BoolP("yes", "y", false, "--yes")
	}

	// simulator flags
	{
		simulatorCmd.PersistentFlags().BoolP("daemon", "", false, "--daemon")
	}

	// exporter flags
	{
		exporterCmd.PersistentFlags().BoolP("customers", "", false, "--customers")
	}

	// scheduler flags
	{
		schedulerCmd.PersistentFlags().BoolP("informings", "", false, "--informings")
		schedulerCmd.PersistentFlags().BoolP("exporter", "", false, "--exporter")
	}

	{
		utils.CompensateDuplicateFeesCmd.PersistentFlags().StringP("start", "", "", "--start=2024-09-04")
		utils.CompensateDuplicateFeesCmd.PersistentFlags().StringP("end", "", "", "--start=2024-09-04")
	}

	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(emitterCmd)
	rootCmd.AddCommand(coreCmd)
	rootCmd.AddCommand(radiusCmd)
	rootCmd.AddCommand(radiusTestCmd)
	rootCmd.AddCommand(feeSchedulerCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(snapshotCmd)
	rootCmd.AddCommand(utilsCmd)
	rootCmd.AddCommand(reporterCmd)
	rootCmd.AddCommand(informingsCmd)
	rootCmd.AddCommand(simulatorCmd)
	rootCmd.AddCommand(exporterCmd)
	rootCmd.AddCommand(dhcpCmd)
	rootCmd.AddCommand(schedulerCmd)
	rootCmd.AddCommand(paymentGatewayCmd)
	rootCmd.AddCommand(utilsCmd)

	utilsCmd.AddCommand(utils.CompensateDuplicateFeesCmd)
}
