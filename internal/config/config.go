package config

import (
	"campaing-comsumer-service/internal/model"
	"context"
	"os"
	"sync"

	"github.com/integralist/go-findroot/find"
	"github.com/lockp111/go-easyzap"
	"github.com/spf13/viper"
)

const (
	ENV_PROFILE_LOCAL = "local"
)

var (
	runOnce sync.Once
	config  model.Config
)

func initConfig() {
	envProfile := os.Getenv("ENV_PROFILE")
	if envProfile == ENV_PROFILE_LOCAL {
		setEnvsByFile()
	}
	viper.AutomaticEnv()
}

func setEnvsByFile() {
	root, _ := find.Repo()

	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath(root.Path + "/build/package/env/local")

	if err := viper.ReadInConfig(); err != nil {
		easyzap.Panic(context.Background(), err, "failed reading config file")
	}
}

func GetConfig() model.Config {
	runOnce.Do(func() {
		initConfig()
		config = model.Config{
			AppName:      viper.GetString("APPLICATION_NAME"),
			ServerPort:   viper.GetString("SERVER_PORT"),
			HealthPort:   viper.GetString("HEALTH_PORT"),
			TimeLocation: viper.GetString("TIME_LOCATION"),
			Database: model.DatabaseConfig{
				Host:     viper.GetString("DATABASE_HOST"),
				Username: viper.GetString("DATABASE_USERNAME"),
				Database: viper.GetString("DATABASE_NAME"),
				Port:     viper.GetInt("DATABASE_PORT"),
				Conn: model.Conn{
					Min:      viper.GetInt("DATABASE_CON_MIN"),
					Max:      viper.GetInt("DATABASE_CON_MAX"),
					Lifetime: viper.GetString("DATABASE_CON_LIFETIME"),
					IdleTime: viper.GetString("DATABASE_CON_IDLETIME"),
				},
			},
		}
		setEnvValues()
	})
	return config
}

func setEnvValues() {
	if len(os.Getenv("DB_PASS_CAMPAING_CONSUMER_API")) > 0 {
		config.Database.Password = os.Getenv("DB_PASS_CAMPAING_CONSUMER_API")
	}
}
