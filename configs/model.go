package configs

import (
	"github.com/go-playground/validator"
	"time"
)

type Config struct {
	Server      Server      `validate:"required"`
	Log         Log         `validate:"required"`
	CorsSupport CorsSupport `validate:"required"`
	JWT         JWT         `validate:"required"`
	DB          DB          `validate:"required"`
	Secrets     Secrets     `validate:"required"`
}

func (c Config) Validate() error {
	return validator.New().Struct(c)
}

type Server struct {
	Name     string        `validate:"required"`
	Port     string        `validate:"required"`
	Timeout  time.Duration `validate:"gt=0"`
	TimeZone string        `validate:"required"`
}

type Log struct {
	Level string `validate:"required"`
	Env   string `validate:"required"`
}

type CorsSupport struct {
	AllowedOrigins string `validate:"required"`
	AllowedMethods string `validate:"required"`
	AllowedHeaders string `validate:"required"`
}

type JWT struct {
	Secret                 string        `validate:"required"`
	ExpirationAccessToken  time.Duration `validate:"gt=0"`
	ExpirationRefreshToken time.Duration `validate:"gt=0"`
}

type DB struct {
	Host         string        `validate:"required"`
	Port         string        `validate:"gt=0"`
	Database     string        `validate:"required"`
	Timeout      time.Duration `validate:"gt=0"`
	MaxIdleConns int           `validate:"gt=0"`
	MaxOpenConns int32         `validate:"gt=0"`
	MaxLifetime  time.Duration `validate:"gt=0"`
}

type Secrets struct {
	DbUsername string `validate:"required"`
	DbPassword string `validate:"required"`
}
