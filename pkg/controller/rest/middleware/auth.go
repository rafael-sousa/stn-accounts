// Package middleware contains custom http request middlewares
package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/response"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

type key string

// Keys used to set request context values
const (
	CtxAccountID key = "CtxAccountID"
)

// NewAuthenticated creates a middleware that requires JWT Authorization Token Header
func NewAuthenticated(jwtH *jwt.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			splitToken := strings.Split(authHeader, " ")
			if len(splitToken) != 2 || !strings.EqualFold("Bearer", splitToken[0]) {
				response.WriteErr(w, r, types.NewErr(types.AuthenticationErr, "authorization header is missing or has an invalid format", nil))
				return
			}
			claims, err := jwtH.Parse(splitToken[1])
			if err != nil {
				if appErr, ok := err.(*types.Err); ok {
					response.WriteErr(w, r, appErr)
					return
				}
				response.WriteErr(w, r, types.NewErr(types.AuthenticationErr, "unable to parse the authorization token", err))
				return
			}
			id, err := strconv.ParseInt(claims.Issuer, 10, 64)
			if err != nil {
				response.WriteErr(w, r, types.NewErr(types.AuthenticationErr, "unable to parse the token issuer", err))
				return
			}
			ctx := context.WithValue(r.Context(), CtxAccountID, id)
			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
