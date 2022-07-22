package db

import (
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var G *gorm.DB

type Config struct {
	Driver      string `toml:"driver"`
	DSN         string `toml:"dsn"`
	MaxIdleConn int    `toml:"max_idle_conn"`
	MaxOpenConn int    `toml:"max_open_conn"`
}

func Init(cfg *Config) error {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		CreateBatchSize:        1000,
	})
	if err != nil {
		return err
	}
	G = db
	return nil
}
