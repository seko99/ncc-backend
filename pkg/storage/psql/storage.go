package psqlstorage

import (
	"code.evixo.ru/ncc/ncc-backend/cmd/config"
	"code.evixo.ru/ncc/ncc-backend/pkg/logger"
	models2 "code.evixo.ru/ncc/ncc-backend/pkg/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	dblogger "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Storage struct {
	db            *gorm.DB
	cfg           *config.Config
	log           logger.Logger
	logLevel      dblogger.LogLevel
	appName       string
	autoMigration bool
}

type StorageOption func(storage *Storage)

func (ths *Storage) GetDB() *gorm.DB {
	return ths.db
}

func (ths *Storage) Connect() error {
	newLogger := dblogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		dblogger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  ths.logLevel, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,        // Disable color
		},
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d application_name='%s' sslmode=disable TimeZone=Europe/Moscow",
		ths.cfg.Db.Host,
		ths.cfg.Db.User,
		ths.cfg.Db.Password,
		ths.cfg.Db.Name,
		ths.cfg.Db.Port,
		ths.appName)

	ths.log.Info(
		"Connecting to DSN: host=%s port=%d dbname=%s application_name='%s'",
		ths.cfg.Db.Host,
		ths.cfg.Db.Port,
		ths.cfg.Db.Name,
		ths.appName,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		PreferSimpleProtocol: true,
		DSN:                  dsn,
	}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	sql, err := db.DB()
	ths.log.Debug(
		"max_open_connections=%d max_idle_connections=%d max_life_time=%v",
		ths.cfg.Db.MaxOpenConnections,
		ths.cfg.Db.MaxIdleConnections,
		ths.cfg.Db.MaxLifeTime,
	)
	sql.SetMaxIdleConns(ths.cfg.Db.MaxIdleConnections)
	sql.SetMaxOpenConns(ths.cfg.Db.MaxOpenConnections)
	sql.SetConnMaxLifetime(ths.cfg.Db.MaxLifeTime)

	ths.db = db

	if ths.autoMigration {
		return ths.Migrate()
	}

	return nil
}

func (ths *Storage) Migrate() error {
	ts := time.Now()
	ths.log.Info("Starting DB migrations...")

	res := ths.GetDB().Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if res.Error != nil {
		return res.Error
	}

	initSQL, err := os.ReadFile("init.sql")
	if err != nil {
		return fmt.Errorf("can't open init.sql")
	}

	res = ths.GetDB().Exec(string(initSQL))
	if res.Error != nil {
		return fmt.Errorf("can't exec init.sql: %w", err)
	}

	err = ths.GetDB().AutoMigrate(
		&models2.CityData{},
		&models2.StreetData{},
		&models2.CustomerGroupData{},

		&models2.BankAccount{},
		&models2.ServiceInternetCustomData{},
		&models2.ServiceInternetData{},
		&models2.ServiceIptvData{},
		&models2.ServiceCatvData{},
		&models2.CustomerFlagData{},
		&models2.ContractData{},
		&models2.IpPoolData{},
		&models2.DeviceStateData{},
		&models2.DeviceGroupData{},
		&models2.HardwareModelData{},
		&models2.MapNodeData{},
		&models2.DeviceData{},
		&models2.ServerGroupData{},
		&models2.ServerData{},
		&models2.NasTypeData{},
		&models2.NasAttributeData{},
		&models2.NasTypeAttributeLink{},
		&models2.NasData{},
		&models2.DhcpPoolData{},
		&models2.PonONUState{},
		&models2.PonONUData{},
		&models2.TriggerData{},
		&models2.IssueTypeData{},
		&models2.IssueUrgencyData{},
		&models2.IssueData{},
		&models2.IssueActionData{},
		&models2.IfaceStateData{},
		&models2.IfaceData{},
		&models2.FDBData{},
		&models2.RadiusVendorData{},
		&models2.RadiusAttributeData{},
		&models2.RadiusAttributeLink{},
		&models2.SecUserData{},
		&models2.SecGroupData{},
		&models2.VendorData{},
		&models2.CustomerData{},
		&models2.SessionData{},
		&models2.SessionsLogData{},
		&models2.DhcpServerData{},
		&models2.DhcpBindingData{},
		&models2.InformingData{},
		&models2.InformingConditionData{},
		&models2.InformingTestCustomerData{},
		&models2.InformingLogData{},
		&models2.FeeLogData{},
		&models2.PaymentData{},
		&models2.PaymentTypeData{},
		&models2.PaymentSystemData{},
		&models2.LeaseData{},
		&models2.LeaseLogData{},
		&models2.PingerStatusData{},
		&models2.CustomerContact{},
		&models2.CustomerServiceLink{},
		&models2.ChartData{},
		&models2.ChartParamData{},
		&models2.ChartParamsLink{},
		&models2.MetricData{},
		&models2.ServiceData{},
		&models2.AccountData{},
		&models2.AccountTransactionData{},
		&models2.PreferencesData{},
		&models2.ScoreProductData{},
		&models2.ScoreLogData{},
		&models2.ScorePaymentTypes{},
		&models2.ScoreExchangeLogData{},
	)
	if err != nil {
		return err
	}

	ths.log.Info("Migration complete in %v", time.Since(ts))

	return nil
}

func WithLogLevel(logLevel dblogger.LogLevel) StorageOption {
	return func(storage *Storage) {
		storage.logLevel = logLevel
	}
}

func WithMigration() StorageOption {
	return func(storage *Storage) {
		storage.autoMigration = true
	}
}

func WithAppName(appName string) StorageOption {
	return func(storage *Storage) {
		storage.appName = appName
	}
}

func NewStorage(cfg *config.Config, log logger.Logger, opts ...StorageOption) *Storage {

	storage := &Storage{
		cfg:      cfg,
		log:      log,
		logLevel: dblogger.LogLevel(cfg.Db.LogLevel),
	}

	for _, opt := range opts {
		opt(storage)
	}

	return storage
}
