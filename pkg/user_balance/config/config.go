package config

import (
	"fmt"
)

type Config struct {
	Host       string
	DbPort     int
	UserDb     string
	UderDbPass string
	DbName     string
	DbConn     string
	Port       int
}

func NewUserBalanceConfig() *Config {

	cfg := Config{
		Host:       "localhost",
		DbPort:     5432,
		UserDb:     "postgres",
		UderDbPass: "12265911",
		DbName:     "mytest",
		Port:       8080,
	}
	cfg.DbConn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.DbPort, cfg.UserDb, cfg.UderDbPass, cfg.DbName)
	return &cfg
}
