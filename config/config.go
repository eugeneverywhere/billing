package config

import "github.com/lillilli/logger"

// Config - service configuration
type Config struct {
	DB     DBConfig
	Rabbit RabbitConfig

	Log logger.Params
}

// RabbitConfig - configuration for rabbit connection
type RabbitConfig struct {
	Addr          string
	OutputChannel string
	ErrorChannel  string
	InputChannel  string
}

// DBConfig - db connection params
type DBConfig struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`

	MaxOpenConnections int
}
