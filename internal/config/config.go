package config

import (
	"flag"

	"github.com/caarlos0/env"
)

type Config struct {
	Settings struct {
		RunAddress           string `env:"RUN_ADDRESS"`
		DatabaseURI          string `env:"DATABASE_URI"`
		AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	} `yaml:"settings"`
}

func GetConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.Settings.RunAddress, "a", ":80", "address for server listen")
	flag.StringVar(&cfg.Settings.DatabaseURI, "d", "", "Database URI")
	flag.StringVar(&cfg.Settings.AccrualSystemAddress, "r", "", "accural system address")

	flag.Parse()

	env.Parse(&cfg.Settings)

	return cfg
}
