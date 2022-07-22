package config

import (
	"github.com/Sora233/DDC/db"
)

type Config struct {
	Addr          string    `toml:"addr"`
	DB            db.Config `toml:"db"`
	LogLevel      string    `toml:"log_level"`
	APIPrefix     string    `toml:"api_prefix"`
	EnableSwagger bool      `toml:"enable_swagger"`
}

var Global Config
