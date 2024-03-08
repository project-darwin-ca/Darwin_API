package config

import (
	env "github.com/caarlos0/env/v6"
	"log"
	"sync"
)

type PostgresConfig struct {
	PostgresHost             string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PostgresPort             string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresDatabase         string `env:"POSTGRES_DATABASE" envDefault:"darwin"`
	PostgresUser             string `env:"POSTGRES_USER" envDefault:"postgres"`
	PostgresPassword         string `env:"POSTGRES_PASSWORD" envDefault:"postgres"`
	PostgresConnectionString string
}

type Admin struct {
	AdminUsername string `env:"ADMIN_USERNAME" envDefault:"admin"`
	AdminPassword string `env:"ADMIN_PASSWORD" envDefault:"password"`
}

type Log struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

type ServiceConfig struct {
	ServiceID       string `env:"SERVICE_ID" envDefault:"data-import"`
	Port            int    `env:"PORT" envDefault:"8080"`
	Profile         string `env:"PROFILE" envDefault:"local"`
	ShutdownTimeout int    `env:"SHUTDOWN_TIMEOUT" endDefault:"3"`
	PostgresConfig
	Admin
	Log
}

func (c ServiceConfig) IsLocalProfile() bool {
	return c.Profile == "local"
}

var once sync.Once
var cfg *ServiceConfig

func GetConfig() *ServiceConfig {
	once.Do(func() {
		config := &ServiceConfig{}
		err := env.Parse(config)
		if err != nil {
			log.Fatalf("failed to parse config: %+v", err)
		}

		cfg = config
		cfg.PostgresConnectionString = "postgres://" +
			cfg.PostgresUser + ":" +
			cfg.PostgresPassword + "@" +
			cfg.PostgresHost + ":" +
			cfg.PostgresPort + "/" +
			cfg.PostgresDatabase
	})
	return cfg
}
