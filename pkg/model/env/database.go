// Package env aggregates configuration models that are used by the application.
// Their values come from the current environment
package env

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
)

// DatabaseConfig maintains the database connection settings
type DatabaseConfig struct {
	Port            int    `env:"DB_PORT,default=3306"`
	User            string `env:"DB_USER,default=admin"`
	Password        string `env:"DB_PW,default=admin"`
	Host            string `env:"DB_HOST,default=localhost"`
	Name            string `env:"DB_NAME,default=stn_accounts"`
	Driver          string `env:"DB_DRIVER,default=mysql"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS,default=10"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS,default=10"`
	ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME,default=0"`
	ParseTime       bool   `env:"DB_PARSE_TIME,default=true"`
}

// DataSourceName builds a datasource name by concatenating the props
func (c *DatabaseConfig) DataSourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=%t", c.User, c.Password, c.Host, c.Port, c.Name, c.ParseTime)
}

// NewDatabaseConfig retrives the environment settings related to the Rest API
func NewDatabaseConfig(ctx *context.Context) DatabaseConfig {
	var c DatabaseConfig
	if err := envconfig.Process(*ctx, &c); err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to read the db application environment properties")
	}
	return c
}
