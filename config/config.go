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
	Db     *DbConfig     `json:"db"`
	Server *ServerConfig `json:"server"`
	Crypto *CryptoConfig `json:"crypto"`
	OAuth  *OAuthConfig  `json:"oauth"`
}

type DbConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type ServerConfig struct {
	Address string `json:"address"`
	Debug   bool   `json:"debug"`
}

type CryptoConfig struct {
	Pepper    string `json:"pepper"`
	HasherKey string `json:"hasher_key"`
}

type OAuthConfig struct {
	Id          string `json:"id"`
	Secret      string `json:"secret"`
	AuthUrl     string `json:"auth_url"`
	TokenUrl    string `json:"token_url"`
	RedirectUrl string `json:"redirect_url"`
}

func LoadConfig(requireFile bool) *Config {
	file, err := os.Open("config.json")

	if err != nil {
		if requireFile {
			panic(err)
		}
		log.Println("no config file present")
		log.Println("using default config")

		return DefaultConfig()
	}
	defer file.Close()

	config := &Config{}

	dec := json.NewDecoder(file)

	err = dec.Decode(config)

	if err != nil {
		log.Println("failed to decode config file")
		log.Println("using default config")

		return DefaultConfig()
	}
	log.Println("loaded config from file")

	return config
}

func DefaultConfig() *Config {
	return &Config{
		Db:     DefaultDbConfig(),
		Server: DefaultServerConfig(),
		Crypto: DefaultCryptoConfig(),
		OAuth:  DefaultOAuthConfig(),
	}

}

func DefaultDbConfig() *DbConfig {
	return &DbConfig{
		Host:     "host",
		Port:     0,
		User:     "user",
		Password: "password",
		Name:     "name",
	}
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Address: "address",
		Debug:   true,
	}
}

func DefaultCryptoConfig() *CryptoConfig {
	return &CryptoConfig{
		Pepper:    "6Sk65RHhGW7S4qnVPV7m",
		HasherKey: "yzzmGPkAA9FTmbtzz9jB",
	}
}

func DefaultOAuthConfig() *OAuthConfig {
	return &OAuthConfig{
		Id:          "id",
		Secret:      "secret",
		AuthUrl:     "auth_url",
		TokenUrl:    "token_url",
		RedirectUrl: "redirect_url",
	}
}

func (c *DbConfig) Dialector() gorm.Dialector {
	var dsn string

	if c.Password == "" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)

	return postgres.Open(dsn)
}
