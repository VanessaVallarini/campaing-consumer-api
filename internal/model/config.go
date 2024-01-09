package model

type Config struct {
	AppName      string
	ServerHost   string
	MetaHost     string
	TimeLocation string
	Database     DatabaseConfig `mapstructure:"database"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	Port     int    `mapstructure:"port"`
	Conn     Conn   `mapstructure:"conn"`
}

type Conn struct {
	Min      int    `mapstructure:"min"`
	Max      int    `mapstructure:"max"`
	Lifetime string `mapstructure:"lifetime"`
	IdleTime string `mapstructure:"idletime"`
}
