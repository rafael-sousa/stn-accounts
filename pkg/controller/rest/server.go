// Package rest and its sub-packages are the application REST entry point.
package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/routing"
	"github.com/rafael-sousa/stn-accounts/pkg/model/env"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

type server struct {
	accountServ  *service.Account
	transferServ *service.Transfer
	middlewares  []func(http.Handler) http.Handler
}

// Server exposes the services provived by the application via REST
type Server interface {
	Use(m ...func(http.Handler) http.Handler) *server
	Start(c *env.RestConfig) error
}

// Use appends the given middleware(s) to the server's middleware slice.
func (s *server) Use(m ...func(http.Handler) http.Handler) *server {
	s.middlewares = append(s.middlewares, m...)
	return s
}

// Start kicks off the service by registering the application routes and applying the configured middlewares
func (s *server) Start(cfg *env.RestConfig) error {

	r := chi.NewRouter()
	for _, md := range s.middlewares {
		r.Use(md)
	}
	jwtH := jwt.NewHandler(cfg)
	r.Route("/accounts", routing.Accounts(s.accountServ))
	r.Route("/transfers", routing.Transfers(s.transferServ, jwtH))
	r.Route("/login", routing.Login(s.accountServ, jwtH))
	r.NotFound(routing.NotFound)
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))
	if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), r); err != nil {
		return fmt.Errorf("Failed to start and listen the http server at port %d, %v", cfg.Port, err)
	}
	return nil
}

// NewServer constructs a server with its required dependencies
func NewServer(accountServ *service.Account, transferServ *service.Transfer) Server {
	return &server{
		accountServ:  accountServ,
		transferServ: transferServ,
	}
}
