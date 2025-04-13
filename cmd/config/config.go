package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type DbConfig struct {
	Host               string        `mapstructure:"host" default:"localhost"`
	Port               int           `mapstructure:"port" default:"5432"`
	Name               string        `mapstructure:"name" default:"ncc"`
	User               string        `mapstructure:"user" default:"ncc"`
	Password           string        `mapstructure:"password" default:"ncc"`
	Schema             string        `mapstructure:"schema" default:"ncc"`
	MaxOpenConnections int           `mapstructure:"max_open_connections" default:"10"`
	MaxIdleConnections int           `mapstructure:"max_idle_connections" default:"10"`
	MaxLifeTime        time.Duration `mapstructure:"max_life_time" default:"1m"`
	LogLevel           int           `mapstructure:"log_level" default:"2"`
}

type DhcpConfig struct {
	Listen string `mapstructure:"listen" default:":8080"`
}

type S3Config struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	AccessKeyID     string `mapstructure:"accesskeyid"`
	SecretAccessKey string `mapstructure:"secretaccesskey"`
	UseSSL          bool   `mapstructure:"usessl"`
	BucketName      string `mapstructure:"bucketname" default:"ncc-backend"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type RadiusServerConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Listen  string `mapstructure:"listen"`
}

type RadiusTestNasConfig struct {
	Ip         string `mapstructure:"ip"`
	Identifier string `mapstructure:"identifier"`
}

type RadiusTestConfig struct {
	Interim time.Duration       `mapstructure:"interim"`
	Limit   int                 `mapstructure:"limit"`
	Secret  string              `mapstructure:"secret"`
	Auth    string              `mapstructure:"auth"`
	Acct    string              `mapstructure:"acct"`
	Nas     RadiusTestNasConfig `mapstructure:"nas"`
}

type RadiusWatcherConfig struct {
	Start   time.Duration `mapstructure:"start"`
	Stop    time.Duration `mapstructure:"stop"`
	Interim time.Duration `mapstructure:"interim"`
}

type RadiusConfig struct {
	Secret  string              `mapstructure:"secret"`
	Update  time.Duration       `mapstructure:"update"`
	Auth    RadiusServerConfig  `mapstructure:"auth"`
	Acct    RadiusServerConfig  `mapstructure:"acct"`
	Control RadiusServerConfig  `mapstructure:"control"`
	Test    RadiusTestConfig    `mapstructure:"test"`
	Watcher RadiusWatcherConfig `mapstructure:"watcher"`
}

type QueueConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type WatcherConfig struct {
	Delay time.Duration `mapstructure:"delay"`
}

type TelegramConfig struct {
	Token  string `mapstructure:"token"`
	ChatID int64  `mapstructure:"chatID"`
}

type ReporterConfig struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
}

type ExporterConfig struct {
	Type         string `mapstructure:"type"`
	Host         string `mapstructure:"host"`
	BadsHost     string `mapstructure:"badshost"`
	Path         string `mapstructure:"path"`
	Username     string `mapstructure:"username"`
	BadsUsername string `mapstructure:"badsusername"`
	Password     string `mapstructure:"password"`
	BadsPassword string `mapstructure:"badspassword"`
	Key          string `mapstructure:"key"`
}

type SmsProviderConfig struct {
	URI   string `mapstructure:"uri"`
	Token string `mapstructure:"token"`
}

type InformingsConfig struct {
	SmsProvider SmsProviderConfig `mapstructure:"smsProvider"`
}

type ApiConfig struct {
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

type PaymentGatewayConfig struct {
	ChargeDailyFee bool `mapstructure:"charge_daily_fee" default:"true"`
	ClearCredit    bool `mapstructure:"clear_credit" default:"true"`
}

type Config struct {
	Listen               string               `mapstructure:"listen"`
	JwtSecret            string               `mapstructure:"jwtSecret"`
	Watcher              WatcherConfig        `mapstructure:"watcher"`
	Db                   DbConfig             `mapstructure:"db"`
	Queue                QueueConfig          `mapstructure:"queue"`
	S3                   S3Config             `mapstructure:"s3"`
	Redis                RedisConfig          `mapstructure:"redis"`
	Dhcp                 DhcpConfig           `mapstructure:"dhcp"`
	Radius               RadiusConfig         `mapstructure:"radius"`
	Reporter             ReporterConfig       `mapstructure:"reporter"`
	Exporter             ExporterConfig       `mapstructure:"exporter"`
	Informings           InformingsConfig     `mapstructure:"informings"`
	API                  ApiConfig            `mapstructure:"api"`
	PaymentGatewayConfig PaymentGatewayConfig `mapstructure:"paygw"`
}

func getDataType(data interface{}) (reflect.Type, error) {
	rv := reflect.ValueOf(data)
	if !rv.IsValid() {
		return nil, fmt.Errorf("can't reflect data")
	}

	// преобразуем к интерфейсу
	iface := rv.Interface()

	rv = reflect.ValueOf(iface)
	if !rv.IsValid() {
		return nil, fmt.Errorf("can't reflect interface")
	}

	// возьмем элемент из интерфейса
	elem := rv.Elem()
	if !elem.IsValid() {
		return nil, fmt.Errorf("can't get reflect elem")
	}

	// определяем тип нижележащей структуры внутри слайса
	// либо непосредственно выставляем тип структуры
	if elem.Type().Kind() == reflect.Slice {
		// создадим новый слайс такого же типа с единственным элементом
		newSlice := reflect.MakeSlice(elem.Type(), 1, 1)

		// берем тип элемента слайса
		return newSlice.Index(0).Type(), nil
	} else if elem.Type().Kind() == reflect.Struct {
		return elem.Type(), nil
	} else {
		return nil, fmt.Errorf("data is not a slice nor single struct")
	}
}

func setEnvValues(v reflect.Value, level int, isStruct bool, defaultValue interface{}, prefix string, p ...[]string) error {
	path := []string{}

	if len(p) > 0 {
		path = p[0]
	}

	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer value")
	}

	v = reflect.Indirect(v)
	if !v.IsValid() {
		return errors.New("value is not valid")
	}

	envVarName := strings.ToUpper(strings.Join(path, "_"))
	if len(prefix) > 0 {
		envVarName = prefix + "_" + envVarName
	}
	envVarValue := os.Getenv(envVarName)

	switch v.Kind() {
	case reflect.Struct:
		if isStruct {
			level++
		}

		for i := 0; i < v.NumField(); i++ {
			fName := v.Type().Field(i).Name

			if len(path) <= level {
				path = append(path, fName)
			} else {
				path[level] = fName
			}

			dv := v.Type().Field(i).Tag.Get("default")
			err := setEnvValues(v.Field(i).Addr(), level, true, dv, prefix, path)
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		if len(envVarValue) > 0 {
			values := strings.Split(envVarValue, ",")
			v.Set(reflect.ValueOf(values))
		} else if defaultValue != nil {
			if val, ok := defaultValue.(string); ok {
				values := strings.Split(val, ",")
				v.Set(reflect.ValueOf(values))
			}
		}
	case reflect.String:
		if len(envVarValue) > 0 {
			v.SetString(envVarValue)
		} else if defaultValue != nil {
			if val, ok := defaultValue.(string); ok {
				v.SetString(val)
			}
		}
	case reflect.Int, reflect.Int64:
		var valToSet string
		if len(envVarValue) > 0 {
			valToSet = envVarValue
		} else if defaultValue != nil {
			if val, ok := defaultValue.(string); ok {
				valToSet = val
			}
		}

		if _, ok := v.Interface().(time.Duration); ok {
			d, err := time.ParseDuration(valToSet)
			if err == nil {
				v.Set(reflect.ValueOf(d))
			}
			break
		} else {
			if len(valToSet) > 0 {
				val, err := strconv.Atoi(valToSet)
				if err != nil {
					return err
				}
				v.SetInt(int64(val))
			}
		}
	case reflect.Bool:
		if len(envVarValue) > 0 {
			val, err := strconv.ParseBool(envVarValue)
			if err != nil {
				return fmt.Errorf("parsing bool env: %w", err)
			}
			v.SetBool(val)
		} else if defaultValue != nil {
			if valStr, ok := defaultValue.(string); ok && valStr != "" {
				val, err := strconv.ParseBool(valStr)
				if err != nil {
					return fmt.Errorf("parsing bool env: %w", err)
				}
				v.SetBool(val)
			}
		}
	default:
	}

	return nil
}

func UnmarshalEnv(cfg interface{}, prefix string) error {

	v := reflect.ValueOf(cfg)
	if reflect.Indirect(v).Kind() != reflect.Struct {
		return errors.New("not a struct")
	}
	return setEnvValues(v, 0, false, nil, prefix)
}

func NewConfig() (*Config, error) {
	var cfg = new(Config)

	err := UnmarshalEnv(cfg, "")
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal env: %w", err)
	}

	return cfg, nil
}
