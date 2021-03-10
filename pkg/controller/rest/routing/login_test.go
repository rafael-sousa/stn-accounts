package routing_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/body"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/routing"
	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestRoutingLoginCreate(t *testing.T) {
	tt := []struct {
		name    string
		service func() service.Account
		status  int
		reader  func() (io.Reader, error)
	}{
		{
			name:   "post '/' successfully",
			status: http.StatusOK,
			service: func() service.Account {
				return &testutil.AccountServMock{
					ExpectLogin: func(c context.Context, cpf string, secret string) (*dto.AccountView, error) {
						testutil.AssertEq(t, "cpf", cpf, "00000000000")
						testutil.AssertEq(t, "secret", secret, "pw")
						return testutil.NewAccountView(1, "Lucas", "00000000000", 999, time.Now()), nil
					},
				}
			},
			reader: func() (io.Reader, error) {
				requestBody := body.LoginRequest{
					CPF:    "00000000000",
					Secret: "pw",
				}
				if body, err := json.Marshal(&requestBody); err == nil {
					return bytes.NewBuffer(body), nil
				} else {
					return nil, err
				}

			},
		},
		{
			name:   "post '/' with empty request body",
			status: http.StatusBadRequest,
			service: func() service.Account {
				return &testutil.AccountServMock{}
			},
			reader: func() (io.Reader, error) {
				if body, err := json.Marshal(""); err == nil {
					return bytes.NewBuffer(body), nil
				} else {
					return nil, err
				}
			},
		},
		{
			name:   "post '/' with invalid request data",
			status: http.StatusBadRequest,
			service: func() service.Account {
				return &testutil.AccountServMock{
					ExpectLogin: func(c context.Context, s1, s2 string) (*dto.AccountView, error) {
						return nil, types.NewErr(types.ValidationErr, "ValidationErr", nil)
					},
				}
			},
			reader: func() (io.Reader, error) {
				requestBody := body.LoginRequest{
					CPF:    "...",
					Secret: "",
				}
				if body, err := json.Marshal(&requestBody); err == nil {
					return bytes.NewBuffer(body), nil
				} else {
					return nil, err
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := chi.NewRouter()
			s := tc.service()
			r.Route("/", routing.Login(&s, jwtHandler))

			buffer, err := tc.reader()
			if err != nil {
				t.Fatalf("unable to create request body")
			}
			req, err := http.NewRequest(http.MethodPost, "/", buffer)
			if err != nil {
				t.Fatalf("unable to create testcase request")
			}
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			testutil.AssertEq(t, "status code", tc.status, res.Code)
		})
	}
}
