// Package rest and its sub-packages are the application REST entry point.
package rest

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/routing"
	"github.com/rafael-sousa/stn-accounts/pkg/model/env"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rs/zerolog/log"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server exposes the services provived by the application via REST
type Server interface {
	Use(m ...func(http.Handler) http.Handler) *server
	Start(c *env.RestConfig)
}

type server struct {
	accountSrv  *service.Account
	transferSrv *service.Transfer
	middlewares []func(http.Handler) http.Handler
}

// NewServer constructs a server with its required dependencies
func NewServer(accountSrv *service.Account, transferSrv *service.Transfer) Server {
	return &server{
		accountSrv:  accountSrv,
		transferSrv: transferSrv,
	}
}

// Use appends the given middleware(s) to the server's middleware slice.
func (s *server) Use(m ...func(http.Handler) http.Handler) *server {
	s.middlewares = append(s.middlewares, m...)
	return s
}

// Start kicks off the service by registering the application routes and applying the configured middlewares
func (s *server) Start(cfg *env.RestConfig) {

	router := chi.NewRouter()
	for _, md := range s.middlewares {
		router.Use(md)
	}
	jwtHandler := jwt.NewHandler(cfg)
	router.Route("/accounts", routing.Accounts(s.accountSrv))
	router.Route("/transfers", routing.Transfers(s.transferSrv, jwtHandler))
	router.Route("/login", routing.Login(s.accountSrv, jwtHandler))
	router.NotFound(routing.NotFound)
	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))

	var httpServer http.Server
	httpServer.Addr = ":" + strconv.Itoa(cfg.Port)
	httpServer.Handler = router

	waitShutdown := make(chan int)
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Info().Msg("Interrupt signal received. Shutting down HTTP server")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Error().Msgf("HTTP server Shutdown: %v", err)
		}
		close(waitShutdown)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Msgf("Failed to start and listen the http server at port %d, %v", cfg.Port, err)
	}

	<-waitShutdown
}
