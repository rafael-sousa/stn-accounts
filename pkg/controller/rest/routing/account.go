package routing

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/response"
	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rs/zerolog/log"
)

type accountHandler struct {
	accountServ *service.Account
}

// @Summary Fetches a list of application accounts
// @tags v1
// @ID fetch-account-list
// @Accept  json
// @Produce  json
// @Success 200 {array} dto.AccountView
// @Failure 500 {object} body.JSONError
// @Router /accounts [get]
func (h accountHandler) get(w http.ResponseWriter, r *http.Request) {
	accs, err := (*h.accountServ).Fetch(r.Context())
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	if err = response.WriteSuccess(w, r, accs, nil); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode the accounts into the response")
		response.WriteErr(w, r, err)
	}
}

// @Summary Gets the current account balance specified by the given ID
// @tags v1
// @ID get-account-balance
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} float64
// @Failure 400 {object} body.JSONError
// @Failure 404 {object} body.JSONError
// @Failure 500 {object} body.JSONError
// @Router /accounts/{id}/balance [get]
func (h *accountHandler) getBalance(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to parse the id param from request URL")
		response.WriteErr(w, r, err)
		return
	}
	b, err := (*h.accountServ).GetBalance(r.Context(), id)
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	if err = response.WriteSuccess(w, r, b, nil); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode the account balance into the response")
		response.WriteErr(w, r, err)
	}
}

// @Summary Creates a new account
// @tags v1
// @ID post-account-create
// @Accept  json
// @Produce  json
// @Param req body dto.AccountCreation required "Account Creation Request"
// @Header 201 {string} Location "/accounts/1"
// @Success 201 {object} dto.AccountView
// @Failure 400 {object} body.JSONError
// @Failure 409 {object} body.JSONError
// @Failure 500 {object} body.JSONError
// @Router /accounts [post]
func (h *accountHandler) post(w http.ResponseWriter, r *http.Request) {
	var d dto.AccountCreation
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to decode the account creation from the request body")
		response.WriteErr(w, r, nil)
		return
	}
	e, err := (*h.accountServ).Create(r.Context(), &d)
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	if err = response.WriteSuccess(w, r, e, e.ID); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode the new account into the response")
		response.WriteErr(w, r, err)
	}
}

// Accounts handle the requests related to entity.Account
func Accounts(accountServ *service.Account) func(chi.Router) {
	h := accountHandler{accountServ: accountServ}
	return func(r chi.Router) {
		r.Get("/", h.get)
		r.Post("/", h.post)
		r.Get("/{id:[\\d]+}/balance", h.getBalance)
	}
}
