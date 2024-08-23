package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local" env-required:"true"`
	HTTPServer `yaml:"http_server"`
	DBServer   `yaml:"db_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	UserApp     string        `yaml:"userApp" env-required:"true"`
	PasswordApp string        `yaml:"passwordApp" env-required:"true" env:"HTTP-SERVER_PASSWORD"`
}

type DBServer struct {
	Dbhost      string `yaml:"dbhost" env-required:"true"`
	Dbport      string `yaml:"dbport" env-required:"true"`
	AliasLenght int    `yaml:"aliasLength" env-required:"true"`
	UserDB      string `yaml:"userDB" env-required:"true"`
	PasswordDB  string `yaml:"passwordDB" env-required:"true" env:"DB-SERVER_PASSWORD"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}
	return &config
}
