package config

import (
	"github.com/hongxeob/go-shopping-practice/server/core/config"
	"github.com/hongxeob/go-shopping-practice/server/core/db"
	"github.com/hongxeob/go-shopping-practice/server/core/endpoint"
	"github.com/hongxeob/go-shopping-practice/server/core/kafka"
	"github.com/hongxeob/go-shopping-practice/server/core/server"
	"go.uber.org/fx"
)

type Config struct {
	Server   server.Config   `yaml:"server"`
	DB       db.Config       `yaml:"db"`
	Endpoint endpoint.Config `yaml:"endpoint"`
	Kafka    kafka.Config    `yaml:"kafka"`
}

func LoadConfig() (Config, error) {
	cfg := &Config{}
	err := config.Unmarshal(cfg)
	if err != nil {
		return Config{}, err
	}
	return *cfg, nil
}

var Module = fx.Options(
	fx.Provide(LoadConfig),
	fx.Provide(func(cfg Config) server.Config {
		return cfg.Server
	}),
	fx.Provide(func(cfg Config) db.Config {
		return cfg.DB
	}),
	fx.Provide(func(cfg Config) kafka.Config {
		return cfg.Kafka
	}),
	fx.Provide(func(cfg Config) endpoint.Config {
		return cfg.Endpoint
	}),
)
