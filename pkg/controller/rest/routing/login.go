package routing

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/body"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/response"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rs/zerolog/log"
)

type loginHandler struct {
	accountSrv *service.Account
	jwtHandler *jwt.Handler
}

// @ID post-login
// @tags v1
// @Summary Generates a new authorization token
// @Accept  json
// @Produce  json
// @Param req body body.LoginRequest required "Login Request"
// @Success 200 {object} body.LoginResponse
// @Failure 400 {object} body.JSONError
// @Failure 500 {object} body.JSONError
// @Router /login [post]
func (h *loginHandler) post(w http.ResponseWriter, r *http.Request) {
	requestBody := body.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to decode request body as body.LoginRequest")
		response.WriteErr(w, r, err)
		return
	}
	view, err := (*h.accountSrv).Login(r.Context(), requestBody.CPF, requestBody.Secret)
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	token, claims, err := (*h.jwtHandler).Generate(view.ID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to generate the jwt token")
		response.WriteErr(w, r, err)
		return
	}
	responseBody := body.LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   int(claims.ExpiresAt) - int(claims.IssuedAt),
	}
	if err = response.WriteSuccess(w, r, responseBody, nil); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode body.LoginResponse into response")
		response.WriteErr(w, r, err)
	}
}

// Login exposes the route that grants user authentication
func Login(accountSrv *service.Account, jwtHandler *jwt.Handler) func(chi.Router) {
	h := loginHandler{
		accountSrv: accountSrv,
		jwtHandler: jwtHandler,
	}
	return func(r chi.Router) {
		r.Post("/", h.post)
	}
}
