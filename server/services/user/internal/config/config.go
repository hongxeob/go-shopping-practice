package config

import (
	"github.com/hongxeob/go-shopping-practice/server/core/db"
	"github.com/hongxeob/go-shopping-practice/server/core/endpoint"
	"github.com/hongxeob/go-shopping-practice/server/core/server"
)

type Config struct {
	Server   server.Config   `yaml:"server"`
	DB       db.Config       `yaml:"db"`
	Endpoint endpoint.Config `yaml:"endpoint"`
}
