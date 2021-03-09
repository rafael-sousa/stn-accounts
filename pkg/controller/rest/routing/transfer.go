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

type transferHandler struct {
	transferSrv *service.Transfer
}

// @ID get-transfer
// @tags v1
// @Summary Gets the list of tranfers for the current authenticated user
// @Accept json
// @Produce json
// @Success 200 {array} dto.TransferView
// @Failure 400 {object} body.JSONError
// @Failure 404 {object} body.JSONError
// @Failure 500 {object} body.JSONError
// @Router /transfers [get]
// @Security ApiKeyAuth
func (h *transferHandler) get(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.CtxAccountID).(int64)
	if !ok {
		response.WriteErr(w, r, types.NewErr(types.InternalErr, "unable to get account id from request context", nil))
		return
	}
	transfers, err := (*h.transferSrv).Fetch(r.Context(), id)
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	if err = response.WriteSuccess(w, r, transfers, nil); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode the transfers into the response")
		response.WriteErr(w, r, err)
	}
}

// @ID post-transfer
// @tags v1
// @Summary Creates a new transfer
// @Accept  json
// @Produce  json
// @Param req body dto.TransferCreation required "Transfer Creation Request"
// @Header 201 {string} Location "/transfers/1"
// @Success 201 {object} dto.TransferView
// @Failure 400 {object} body.JSONError
// @Failure 404 {object} body.JSONError
// @Failure 500 {object} body.JSONError
// @Router /transfers [post]
// @Security ApiKeyAuth
func (h *transferHandler) post(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(middleware.CtxAccountID).(int64)
	if !ok {
		response.WriteErr(w, r, types.NewErr(types.InternalErr, "unable to get account id from request context", nil))
		return
	}
	var transferCreation dto.TransferCreation
	err := json.NewDecoder(r.Body).Decode(&transferCreation)
	if err != nil {
		log.Error().Caller().Err(err).Msg("unable to decode request body as transfer creation")
		response.WriteErr(w, r, err)
		return
	}

	view, err := (*h.transferSrv).Create(r.Context(), id, &transferCreation)
	if err != nil {
		response.WriteErr(w, r, err)
		return
	}
	if err = response.WriteSuccess(w, r, view, view.ID); err != nil {
		log.Error().Caller().Err(err).Msg("unable to encode transfer into response")
		response.WriteErr(w, r, err)
	}
}

// Transfers handles the requests related to entity.Transfer
func Transfers(transferSrv *service.Transfer, jwtHandler *jwt.Handler) func(chi.Router) {
	h := transferHandler{transferSrv: transferSrv}
	return func(r chi.Router) {
		r.Use(middleware.NewAuthenticated(jwtHandler))
		r.Get("/", h.get)
		r.Post("/", h.post)
	}
}
