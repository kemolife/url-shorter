package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default: "local"`
	Storage    `yaml:"storage"`
	HttpServer `yaml:"http_server"`
	GitHub     `yaml:"github"`
	Auth       `yaml:"auth"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default: "localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default: "4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default: "4s"`
}

type Storage struct {
	Type   string            `yaml:"type" env-default: "sqlite"`
	Config map[string]string `yaml:"config"`
}

type GitHub struct {
	ClientId     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

type Auth struct {
	JWTSecretKey     string `yaml:"jwt_secret_key"`
	AllowedGitHubOrg string `yaml:"git_hub_org"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatalf("CONFIG_PATH is not set")
		}

		fmt.Println(configPath)
		//check if file exist
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file not exist")
		}

		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("cannot read config %s", err)
		}
	})

	return cfg
}
