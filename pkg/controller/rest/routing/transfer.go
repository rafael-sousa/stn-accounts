package routing

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/middleware"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/response"
	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rs/zerolog/log"
)

type router struct {
	transferServ *service.Transfer
}

func (h *router) get(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.CtxAccountID).(int64)
	if !ok {
		response.WriteErr(w, r, types.NewErr(types.InternalErr, "unable to get account id from the current request context", nil))
		return
	}
	transfers, err := (*h.transferServ).Fetch(r.Context(), id)
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	if err = response.WriteSuccess(w, r, transfers, nil); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode the transfers into the response")
		response.WriteErr(w, r, err)
	}
}

func (h *router) post(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.CtxAccountID).(int64)
	if !ok {
		response.WriteErr(w, r, types.NewErr(types.InternalErr, "unable to get account id from the current request context", nil))
		return
	}
	var d dto.TransferCreation
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to decode the request payload as a transfer creation")
		response.WriteErr(w, r, err)
		return
	}

	t, err := (*h.transferServ).Create(r.Context(), id, &d)
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	if err = response.WriteSuccess(w, r, t, t.ID); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode the transfer into the response")
		response.WriteErr(w, r, err)
	}
}

// Transfers handles the requests related to entity.Transfer
func Transfers(transferServ *service.Transfer, jwtH *jwt.Handler) func(chi.Router) {
	h := router{transferServ: transferServ}
	return func(r chi.Router) {
		r.Use(middleware.NewAuthenticated(jwtH))
		r.Get("/", h.get)
		r.Post("/", h.post)
	}
}
