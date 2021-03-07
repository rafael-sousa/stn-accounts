package env

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

// RestConfig maintains the configuration for the Rest API
type RestConfig struct {
	Port            int    `env:"PORT,default=3000"`
	Secret          []byte `env:"JWT_SECRET,default=rest-app@@secret"`
	TokenExpTimeout int    `env:"JWT_EXP_TIMEOUT,default=30"`
}

// NewRestConfig retrives the environment settings related to the Rest API
func NewRestConfig(ctx *context.Context) RestConfig {
	var c RestConfig
	if err := envconfig.Process(*ctx, &c); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to read the rest application environment properties")
	}
	return c
}
