package config

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/joomcode/errorx"
	"github.com/lockp111/go-easyzap"
	"github.com/spf13/viper"
)

type Config struct {
	AppName      string
	ServerPort   string
	HealthPort   string
	TimeLocation string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Database string
	Password string
	Conn     Conn
}

type Conn struct {
	Min      int    `mapstructure:"min"`
	Max      int    `mapstructure:"max"`
	Lifetime string `mapstructure:"lifetime"`
	IdleTime string `mapstructure:"idletime"`
}

type AwsConfig struct {
	Url           string
	Region        string
	QueueCampaing string
	Credentials   aws.AnonymousCredentials
}

var config = viper.New()

func Init() {
	config.AddConfigPath("internal/config/")
	config.SetConfigName("configuration")
	config.SetConfigType("yml")

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			err := errorx.Decorate(err, "Error reading config file: %s", err)
			easyzap.Fatal(context.Background(), err, "Unable to keep the service without config file")
		}
	}

	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	config.AutomaticEnv()
}

func GetDatabaseConfig() DatabaseConfig {
	if err := config.BindEnv("dataBase.password", "DB_PASS_CAMPAING_CONSUMER_API"); err != nil {
		errorx.Panic(err)
	}

	databaseConfig := DatabaseConfig{
		Host:     config.GetString("database.host"),
		Port:     config.GetInt("database.port"),
		Username: config.GetString("database.username"),
		Database: config.GetString("database.database"),
		Password: config.GetString("database.password"),
		Conn: Conn{
			Min:      config.GetInt("database.conn.min"),
			Max:      config.GetInt("database.conn.max"),
			Lifetime: config.GetString("database.conn.lifetime"),
			IdleTime: config.GetString("database.conn.idletime"),
		},
	}

	return databaseConfig
}

func GetAwsConfig() AwsConfig {
	return AwsConfig{
		Url:           config.GetString("aws.url"),
		Region:        config.GetString("aws.region"),
		QueueCampaing: config.GetString("aws.sqs.queue"),
	}
}

func GetConfig() Config {
	return Config{
		AppName:      config.GetString("app.name"),
		ServerPort:   config.GetString("server.port"),
		HealthPort:   config.GetString("health.port"),
		TimeLocation: config.GetString("timeLocation"),
	}
}
