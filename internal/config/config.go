package config

import (
	"campaing-comsumer-service/internal/model"
	"context"
	"os"
	"strings"
	"sync"

	"github.com/lockp111/go-easyzap"
	"github.com/spf13/viper"
)

var (
	runOnce sync.Once
	config  model.Config
)

func GetConfig() model.Config {
	runOnce.Do(func() {
		viperConfig := initConfig()

		config = model.Config{
			AppName:      viperConfig.GetString("app.name"),
			ServerHost:   viperConfig.GetString("server.host"),
			MetaHost:     viperConfig.GetString("meta.host"),
			TimeLocation: viperConfig.GetString("time-location"),
			Database: model.DatabaseConfig{
				Host:     viperConfig.GetString("database.host"),
				Username: viperConfig.GetString("database.username"),
				Database: viperConfig.GetString("database.database"),
				Port:     viperConfig.GetInt("database.port"),
				Conn: model.Conn{
					Min:      viperConfig.GetInt("database.conn.min"),
					Max:      viperConfig.GetInt("database.conn.max"),
					Lifetime: viperConfig.GetString("database.conn.lifetime"),
					IdleTime: viperConfig.GetString("database.conn.idletime"),
				},
			},
		}
		setEnvValues()
	})
	return config
}

func initConfig() viper.Viper {
	config := viper.New()

	initDefaults(config)
	environment := os.Getenv("ENV_PROFILE")

	config.SetConfigType("yaml")
	config.AddConfigPath("./internal/config")
	config.SetConfigName(environment)

	err := config.MergeInConfig()
	if err != nil {
		easyzap.Fatal(context.TODO(), err, "error reading application config file.")
	}

	return *config
}

func initDefaults(config *viper.Viper) {
	config.SetDefault("app.name", "campaing-consumer-api")
	config.SetDefault("server.host", "0.0.0.0:8080")
	config.SetDefault("meta.host", "0.0.0.0:8081")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
}

func setEnvValues() {
	if len(os.Getenv("DB_PASS_CAMPAING_CONSUMER_API")) > 0 {
		config.Database.Password = os.Getenv("DB_PASS_CAMPAING_CONSUMER_API")
	}
}
