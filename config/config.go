package config

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type Config struct {
	Postgres *PostgresConfig `json:"postgres"`
	Server   *ServerConfig   `json:"server"`
	Crypto   *CryptoConfig   `json:"crypto"`
}

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"db_name"`
}

type ServerConfig struct {
	Address string `json:"address"`
	IsDebug bool   `json:"is_debug"`
}

type CryptoConfig struct {
	Pepper    string `json:"pepper"`
	HasherKey string `json:"hasher_key"`
}

func LoadConfig(requireFile bool) *Config {
	file, err := os.Open("config.json")

	if err != nil {
		if requireFile {
			panic(err)
		}
		log.Println("no config.json present")
		log.Println("using default config")

		return DefaultConfig()
	}
	defer file.Close()

	config := &Config{}

	dec := json.NewDecoder(file)

	err = dec.Decode(config)

	if err != nil {
		log.Println("failed to unmarshal config.json")
		log.Println("using default config")

		return DefaultConfig()
	}
	log.Println("loaded config.json")

	return config
}

func DefaultConfig() *Config {
	return &Config{
		Postgres: DefaultPostgresConfig(),
		Server:   DefaultServerConfig(),
		Crypto:   DefaultCryptoConfig(),
	}

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

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Address: "localhost:8080",
		IsDebug: true,
	}
}

func DefaultCryptoConfig() *CryptoConfig {
	return &CryptoConfig{
		Pepper:    "6Sk65RHhGW7S4qnVPV7m",
		HasherKey: "yzzmGPkAA9FTmbtzz9jB",
	}
}

func (c *PostgresConfig) Dialector() gorm.Dialector {
	var dsn string

	if c.Password == "" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Dbname)
	}
	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Dbname)

	return postgres.Open(dsn)
}
