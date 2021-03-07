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
	accountServ *service.Account
	jwtH        *jwt.Handler
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
	b := body.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to decode the request payload as a login creation")
		response.WriteErr(w, r, err)
		return
	}
	view, err := (*h.accountServ).Login(r.Context(), b.CPF, b.Secret)
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	token, claims, err := (*h.jwtH).Generate(view.ID)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to generate the jwt token")
		response.WriteErr(w, r, err)
		return
	}
	p := body.LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   int(claims.ExpiresAt) - int(claims.IssuedAt),
	}
	if err = response.WriteSuccess(w, r, p, nil); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode the login info into the response")
		response.WriteErr(w, r, err)
	}
}

// Login exposes the route that grants user authentication
func Login(accountServ *service.Account, jwtH *jwt.Handler) func(chi.Router) {
	h := loginHandler{
		accountServ: accountServ,
		jwtH:        jwtH,
	}
	return func(r chi.Router) {
		r.Post("/", h.post)
	}
}
