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

// Generate creates a new JWT token with the given id
func (h *Handler) Generate(id int64) (string, *jwtgo.StandardClaims, error) {
	c := jwtgo.StandardClaims{
		ExpiresAt: time.Now().Add(h.expTimeout * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    strconv.FormatInt(id, 10),
	}
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, &c)
	tokenStr, err := token.SignedString(h.secret)
	return tokenStr, &c, err
}

// Parse receives a token string, parses it, validates its content, and returns a token claims
func (h *Handler) Parse(tokenHeader string) (*jwtgo.StandardClaims, error) {
	c := &jwtgo.StandardClaims{}
	token, err := jwtgo.ParseWithClaims(tokenHeader, c, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, types.NewErr(types.AuthenticationErr, "invalid jwt token signature", nil)
		}
		if _, ok := token.Claims.(*jwtgo.StandardClaims); !ok {
			return nil, types.NewErr(types.AuthenticationErr, "unexpected token content", nil)
		}
		return h.secret, nil
	})
	if ve, ok := err.(*jwtgo.ValidationError); ok {
		if ve.Errors&jwtgo.ValidationErrorMalformed != 0 {
			return nil, types.NewErr(types.AuthenticationErr, "malformed jwt token", nil)
		} else if ve.Errors&(jwtgo.ValidationErrorExpired|jwtgo.ValidationErrorNotValidYet) != 0 {
			return nil, types.NewErr(types.AuthenticationErr, "expired or premature jwt token", nil)
		}
	}
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return c, nil
	}
	return nil, types.NewErr(types.AuthenticationErr, "unexpected token format", nil)
}

// NewHandler creates a new JWT Handler
func NewHandler(c *env.RestConfig) *Handler {
	return &Handler{
		secret:     c.Secret,
		expTimeout: time.Duration(c.TokenExpTimeout),
	}
}
