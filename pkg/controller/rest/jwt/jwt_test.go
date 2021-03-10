package jwt_test

import (
	"strconv"
	"testing"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/model/env"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestGenerate(t *testing.T) {
	tt := []struct {
		name   string
		config *env.RestConfig
		input  int64
	}{
		{
			name: "generate jwt token successfully",
			config: &env.RestConfig{
				TokenExpTimeout: 30,
				Secret:          []byte("secret"),
			},
			input: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h := jwt.NewHandler(tc.config)
			if token, claims, err := h.Generate(tc.input); err == nil {
				if len(token) == 0 {
					t.Error("expected not empty token")
				}
				testutil.AssertEq(t, "token issuer", claims.Issuer, strconv.FormatInt(tc.input, 10))
				currentTimeout := claims.ExpiresAt - claims.IssuedAt
				expectedTimeout := time.Duration(tc.config.TokenExpTimeout) * time.Minute

				testutil.AssertEq(t, "token timeout", int64(expectedTimeout.Seconds()), currentTimeout)
			} else {
				t.Error(err)
			}

		})
	}
}

func TestParse(t *testing.T) {
	generate := func(claims *jwtgo.StandardClaims) (string, error) {
		token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
		return token.SignedString([]byte("secret"))
	}

	tt := []struct {
		name      string
		input     *jwtgo.StandardClaims
		config    *env.RestConfig
		assertErr func(*testing.T, error)
	}{
		{
			name: "parse jwt token successfully",
			input: &jwtgo.StandardClaims{
				ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "1",
			},
			config: &env.RestConfig{
				Secret: []byte("secret"),
			},
			assertErr: func(t *testing.T, e error) { t.Error(e) },
		},
		{
			name: "parse jwt token expired",
			input: &jwtgo.StandardClaims{
				ExpiresAt: time.Now().Add(-time.Minute).Unix(),
				IssuedAt:  time.Now().Add(-time.Hour).Unix(),
				Issuer:    "1",
			},
			config: &env.RestConfig{
				Secret: []byte("secret"),
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.AuthenticationErr, err, "expired or premature jwt token")
			},
		},
		{
			name: "parse jwt token ...",
			input: &jwtgo.StandardClaims{
				ExpiresAt: time.Now().Add(time.Minute).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "1",
			},
			config: &env.RestConfig{
				Secret: []byte("oba"),
			},
			assertErr: func(t *testing.T, err error) {
				testutil.AssertCustomErr(t, types.AuthenticationErr, err, "invalid token signature")
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			h := jwt.NewHandler(tc.config)
			token, err := generate(tc.input)
			if err != nil {
				t.Error(err)
				return
			}
			if claims, err := h.Parse(token); err == nil {
				testutil.AssertEq(t, "token issuer", tc.input.Issuer, claims.Issuer)
			} else {
				tc.assertErr(t, err)
			}
		})
	}
}
