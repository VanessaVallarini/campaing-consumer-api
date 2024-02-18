package config

import (
	"context"
	"fmt"
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
	PostgresDriver  string
	User            string
	Host            string
	Port            int
	Password        string
	DbName          string
	Conn            int
	DatabaseConnStr string
	ConnMax         int
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

	setConfigDefaults()

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

func setConfigDefaults() {
	config.SetDefault("app.name", "campaing-consumer-api")
	config.SetDefault("server.port", "0.0.0.0:8080")
	config.SetDefault("health.port", "0.0.0.0:8081")
	config.SetDefault("timeLocation", "America/Sao_Paulo")
	config.SetDefault("dataBase.driver", "postgres")
	config.SetDefault("dataBase.host", "pg-campaing-consumer-api.sandbox.com")
	config.SetDefault("dataBase.port", "5432")
	config.SetDefault("dataBase.user", "ads-campaing-app")
	config.SetDefault("dataBase.name", "ads-campaing-db")
	config.SetDefault("dataBase.connMax", "20")
	config.SetDefault("aws.url", "https://sqs.us-east-1.amazonaws.com")
	config.SetDefault("aws.region", "us-east-1")
	config.SetDefault("aws.sqs.queue", "https://sqs.us-east-1.amazonaws.com/queue_campaing")
}

func GetDatabaseConfig() DatabaseConfig {
	if err := config.BindEnv("dataBase.password", "DB_PASS_CAMPAING_CONSUMER_API"); err != nil {
		errorx.Panic(err)
	}

	databaseConfig := DatabaseConfig{
		PostgresDriver: config.GetString("dataBase.driver"),
		User:           config.GetString("dataBase.user"),
		Host:           config.GetString("dataBase.host"),
		Port:           config.GetInt("dataBase.port"),
		DbName:         config.GetString("dataBase.name"),
		ConnMax:        config.GetInt("dataBase.connMax"),
		Password:       config.GetString("dataBase.password"),
	}

	databaseConfig.DatabaseConnStr = buildDatabaseConnString(databaseConfig)

	return databaseConfig
}

func buildDatabaseConnString(dbCfg DatabaseConfig) string {
	connectionDSN := fmt.Sprintf("user=%s host=%s port=%v  "+
		"password=%s dbname=%s connect_timeout=%v sslmode=disable",
		dbCfg.User,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Password,
		dbCfg.DbName,
		dbCfg.Conn,
	)

	return connectionDSN
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
