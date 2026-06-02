package core_jwt

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Expiration time.Duration `envconfig:"EXPIRATION" required:"true"`
	Secret     string        `envconfig:"SECRET" required:"true"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("JWT", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get JWT config: %w", err)
		panic(err)
	}

	return config
}
