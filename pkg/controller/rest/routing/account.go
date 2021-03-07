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

func (h *accountHandler) getBalance(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to parse the id from param from the request URL")
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

func (h *accountHandler) post(w http.ResponseWriter, r *http.Request) {
	var d dto.AccountCreation
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to decode the account creation from the request body")
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
