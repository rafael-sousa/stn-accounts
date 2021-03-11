// Package jwt contains structs that handles JWT token encoding and parsing
package jwt

import (
	"strconv"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/rafael-sousa/stn-accounts/pkg/model/env"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

// Handler exposes JWT related functions
type Handler struct {
	secret     []byte
	expTimeout time.Duration
}

// NewHandler creates a new JWT Handler
func NewHandler(config *env.RestConfig) *Handler {
	return &Handler{
		secret:     config.Secret,
		expTimeout: time.Duration(config.TokenExpTimeout),
	}
}

// Generate creates a new JWT token with the given id
func (h *Handler) Generate(id int64) (string, *jwtgo.StandardClaims, error) {
	claims := jwtgo.StandardClaims{
		ExpiresAt: time.Now().Add(h.expTimeout * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    strconv.FormatInt(id, 10),
	}
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, &claims)
	signedString, err := token.SignedString(h.secret)
	return signedString, &claims, err
}

// Parse receives a token string, parses it, validates its content, and returns a token claims
func (h *Handler) Parse(tokenString string) (*jwtgo.StandardClaims, error) {
	claims := &jwtgo.StandardClaims{}
	token, err := jwtgo.ParseWithClaims(tokenString, claims, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, types.NewErr(types.AuthenticationErr, "invalid jwt token signature", nil)
		}
		if _, ok := token.Claims.(*jwtgo.StandardClaims); !ok {
			return nil, types.NewErr(types.AuthenticationErr, "unexpected token content", nil)
		}
		return h.secret, nil
	})
	if ve, ok := err.(*jwtgo.ValidationError); ok {
		switch {
		case ve.Errors&jwtgo.ValidationErrorMalformed != 0:
			return nil, types.NewErr(types.AuthenticationErr, "malformed jwt token", nil)
		case ve.Errors&(jwtgo.ValidationErrorExpired|jwtgo.ValidationErrorNotValidYet) != 0:
			return nil, types.NewErr(types.AuthenticationErr, "expired or premature jwt token", nil)
		case ve.Errors&jwtgo.ValidationErrorSignatureInvalid != 0:
			return nil, types.NewErr(types.AuthenticationErr, "invalid token signature", nil)
		}
	}
	if err != nil || !token.Valid {
		return nil, types.NewErr(types.AuthenticationErr, "unexpected token format", err)
	}
	return claims, nil
}
