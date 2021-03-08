package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-chi/chi/middleware"
	_ "github.com/rafael-sousa/stn-accounts/docs"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest"
	"github.com/rafael-sousa/stn-accounts/pkg/model/env"
	"github.com/rafael-sousa/stn-accounts/pkg/repository"
	"github.com/rafael-sousa/stn-accounts/pkg/repository/mysql"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
)

// @title Account REST API
// @version 1.0
// @description Application that exposes a REST API.
//
// @contact.name Rafael S.
// @contact.url github.com/rafael-sousa
// @contact.email rafaelsj7@gmail.com
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
// @license.name MIT
// @license.url https://github.com/rafael-sousa/stn-accounts/blob/main/LICENSE
func main() {
	// Parse the application configuration
	ctx := context.Background()
	dbConfig := env.NewDatabaseConfig(&ctx)
	restConfig := env.NewRestConfig(&ctx)

	// Set up a database connection pool
	db, err := sql.Open(dbConfig.Driver, dbConfig.DataSourceName())
	if err != nil {
		log.Fatal().
			Caller().
			Err(err).
			Str("driver", dbConfig.Driver).
			Str("db_name", dbConfig.Name).
			Str("db_host", dbConfig.Host).
			Int("db_port", dbConfig.Port).
			Msg("Unable to open database connection pool")
	}
	defer db.Close()

	db.SetConnMaxLifetime(time.Minute * time.Duration(dbConfig.ConnMaxLifetime))
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)

	// Initializes the application dependency tree
	txr := repository.NewTxr(db)
	accountRepo := mysql.NewAccount(&txr)
	transferRepo := mysql.NewTransfer(&txr)
	accountServ := service.NewAccount(&txr, &accountRepo)
	transferServ := service.NewTransfer(&txr, &transferRepo, &accountRepo)
	server := rest.NewServer(&accountServ, &transferServ)

	server.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

	server.Start(&restConfig)
}
