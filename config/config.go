package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"db_name"`
}

func DefaultPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "Abcde12345!",
		Dbname:   "lenslocked_dev",
	}
}

func (c *PostgresConfig) dsn() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Dbname)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Dbname)
}

func (c *PostgresConfig) Dialector() gorm.Dialector {
	return postgres.Open(c.dsn())
}

type ServerConfig struct {
	Address string `json:"address"`
	IsDebug bool   `json:"is_debug"`
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Address: "localhost:8080",
		IsDebug: true,
	}
}

type CryptoConfig struct {
	Pepper    string `json:"pepper"`
	HasherKey string `json:"hasher_key"`
}

func DefaultCryptoConfig() *CryptoConfig {
	return &CryptoConfig{
		Pepper:    "6Sk65RHhGW7S4qnVPV7m",
		HasherKey: "yzzmGPkAA9FTmbtzz9jB",
	}
}
