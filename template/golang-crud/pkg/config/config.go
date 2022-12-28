package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/gin-contrib/requestid"
	"github.com/rs/zerolog"
)

type Config struct {
	Logger
	MongoDB
	Pagination
	Routes
	Server
}

type Logger struct {
	Level           zerolog.Level          `env:"LOGGER_LEVEL,required,notEmpty" envDefault:"info"`
	RequestIDHeader requestid.HeaderStrKey `env:"LOGGER_REQUEST_ID_HEADER,required,notEmpty" envDefault:"x-correlation-id"`
}

type MongoDB struct {
	URL    string `env:"MONGODB_URL,required,notEmpty"`
	DBName string `env:"MONGODB_DB_NAME,required,notEmpty" envDefault:"openfaas-golang-crud"`
}

type Pagination struct {
	Limit    int64 `env:"PAGINATION_LIMIT,required,notEmpty" envDefault:"25"`
	MaxLimit int64 `env:"PAGINATION_MAX_LIMIT,required,notEmpty" envDefault:"100"`
}

type Routes struct {
	Create  bool `env:"CREATE_ONE_ROUTE_ENABLED" envDefault:"true"`
	Delete  bool `env:"DELETE_ONE_ROUTE_ENABLED" envDefault:"true"`
	GetMany bool `env:"GET_MANY_ROUTE_ENABLED" envDefault:"true"`
	GetOne  bool `env:"GET_ONE_ROUTE_ENABLED" envDefault:"true"`
	Update  bool `env:"UPDATE_ONE_ROUTE_ENABLED" envDefault:"true"`
}

type Server struct {
	RoutePrefix string `env:"ROUTE_PREFIX,required,notEmpty" envDefault:"/crud"`
	Port        int    `env:"http_port,required,notEmpty" envDefault:"3000"` // Use http_port for compatibility with OpenFaaS watchdog
}

func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
