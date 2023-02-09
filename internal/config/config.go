package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env"
)

type Config struct {
	Settings struct {
		RunAddress           string        `env:"RUN_ADDRESS"`
		DatabaseURI          string        `env:"DATABASE_URI"`
		AccrualSystemAddress string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
		AccrualScanInterval  time.Duration `env:"ACCRUAL_SCAN_INTERVAL"`
	}
}

func GetConfig() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.Settings.RunAddress, "a", ":80", "address for server listen")
	flag.StringVar(&cfg.Settings.DatabaseURI, "d", "", "Database URI")
	flag.StringVar(&cfg.Settings.AccrualSystemAddress, "r", "", "accural system address")
	flag.DurationVar(&cfg.Settings.AccrualScanInterval, "s", time.Millisecond*150, "accural scan interval")

	flag.Parse()

	env.Parse(&cfg.Settings)

	return cfg
}
