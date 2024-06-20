package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env          string         `yaml:"env"`
	StartTimeout time.Duration  `yaml:"start_timeout"`
	GRPC         GRPCConfig     `yaml:"grpc"`
	Redis        RedisConfig    `yaml:"redis"`
	Postgres     PostgresConfig `yaml:"postgres"`
}

const (
	EnvLocal = "local"
	EnvProd  = "prod"
	EnvDev   = "dev"
)

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type RedisConfig struct {
	Host     string        `yaml:"host"`
	Port     string        `yaml:"port"`
	DB       int           `yaml:"db"`
	Password string        `yaml:"password"`
	Timeout  time.Duration `yaml:"timeout"`
}

type PostgresConfig struct {
	Host               string        `yaml:"host"`
	Port               string        `yaml:"port"`
	User               string        `yaml:"user"`
	Password           string        `yaml:"password"`
	DBName             string        `yaml:"db"`
	SSLMode            string        `yaml:"ssl_mode"`
	MaxConnections     int           `yaml:"max_connections"`
	MaxIdleConnections int           `yaml:"max_idle_connections"`
	Timeout            time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Fatalf("config file is not exists: %s", cfgPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s:%s", cfgPath, err)
	}

	return &cfg
}
