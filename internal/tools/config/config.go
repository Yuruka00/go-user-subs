package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	PgPort     string `env:"POSTGRES_PORT,required" envDefault:"5432"`
	PgHost     string `env:"POSTGRES_HOST,required"`
	PgDB       string `env:"POSTGRES_DB,required"`
	PgUser     string `env:"POSTGRES_USER,required"`
	PgPassword string `env:"POSTGRES_PASSWORD,required"`
	AppPort    string `env:"APP_PORT,required"`
}

func Load() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("config parse error: %w", err)
	}

	return &cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		c.PgHost, c.PgUser, c.PgPassword, c.PgDB, c.PgPort,
	)
}
