package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/middleware"
	"github.com/rafael-sousa/stn-accounts/pkg/model/env"
)

func TestAuthorizationRequest(t *testing.T) {
	generate := func(claims *jwtgo.StandardClaims) (string, error) {
		token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
		return token.SignedString([]byte("secret"))
	}

	jwtHandler := jwt.NewHandler(&env.RestConfig{
		TokenExpTimeout: 30,
		Secret:          []byte("secret"),
	})

	tt := []struct {
		name           string
		claims         *jwtgo.StandardClaims
		assertResponse func(*testing.T, *http.Response)
	}{
		{
			name: "intercept request with valid auth header",
			claims: &jwtgo.StandardClaims{
				ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
				IssuedAt:  time.Now().Unix(),
				Issuer:    "1",
			},
			assertResponse: func(t *testing.T, r *http.Response) {
				if r.StatusCode != http.StatusOK {
					t.Errorf("expected status code '%d' but got '%d'", http.StatusOK, r.StatusCode)
				}
			},
		},
		{
			name: "intercept request with no auth header",
			assertResponse: func(t *testing.T, r *http.Response) {
				if r.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected status code '%d' but got '%d'", http.StatusUnauthorized, r.StatusCode)
				}
			},
		},
		{
			name: "intercept request with expired auth header",
			claims: &jwtgo.StandardClaims{
				ExpiresAt: time.Now().Add(-time.Minute).Unix(),
				IssuedAt:  time.Now().Add(-time.Hour).Unix(),
				Issuer:    "1",
			},
			assertResponse: func(t *testing.T, r *http.Response) {
				if r.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected status code '%d' but got '%d'", http.StatusUnauthorized, r.StatusCode)
				}
			},
		},
		{
			name: "intercept request with invalid issuer",
			claims: &jwtgo.StandardClaims{
				Issuer: "foo",
			},
			assertResponse: func(t *testing.T, r *http.Response) {
				if r.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected status code '%d' but got '%d'", http.StatusUnauthorized, r.StatusCode)
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			m := middleware.NewAuthenticated(jwtHandler)

			handler := m(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if id, ok := r.Context().Value(middleware.CtxAccountID).(int64); ok {
					if tc.claims.Issuer != strconv.FormatInt(id, 10) {
						t.Errorf("expected issuer id equal to '%s' but got '%v'", tc.claims.Issuer, id)
					}
				} else {
					t.Errorf("unabled to retrieve issuer id from request")
				}
			}))

			request, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Error(err)
				return
			}
			if tc.claims != nil {
				if token, err := generate(tc.claims); err == nil {
					request.Header.Set("Authorization", "Bearer "+token)
				} else {
					t.Error(err)
					return
				}
			}
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			tc.assertResponse(t, response.Result())

		})
	}
}
