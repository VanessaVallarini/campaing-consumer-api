package model

type Config struct {
	AppName      string
	ServerPort   string
	HealthPort   string
	TimeLocation string
	Database     DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Username string
	Password string
	Database string
	Port     int
	Conn     Conn
}

type Conn struct {
	Min      int
	Max      int
	Lifetime string
	IdleTime string
}
