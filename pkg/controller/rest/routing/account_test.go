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
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/routing"
	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestRoutingAccountFetch(t *testing.T) {
	tt := []struct {
		name      string
		service   func() service.Account
		status    int
		path      string
		assertRes func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "get '/' with no result successfully",
			status: http.StatusOK,
			path:   "/",
			service: func() service.Account {
				return &testutil.AccountServMock{
					ExpectFetch: func(c context.Context) ([]dto.AccountView, error) {
						return []dto.AccountView{}, nil
					},
				}
			},
		},
		{
			name:   "get '/' with results successfully",
			status: http.StatusOK,
			path:   "/",
			service: func() service.Account {
				return &testutil.AccountServMock{
					ExpectFetch: func(c context.Context) ([]dto.AccountView, error) {
						return []dto.AccountView{
							*testutil.NewAccountView(1, "Jose", "00000000000", 5, time.Now()),
							*testutil.NewAccountView(2, "Silva", "11111111111", 10, time.Now()),
							*testutil.NewAccountView(3, "Sousa", "22222222222", 20, time.Now()),
						}, nil
					},
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := chi.NewRouter()
			s := tc.service()
			r.Route("/", routing.Accounts(&s))

			req, err := http.NewRequest(http.MethodGet, tc.path, nil)
			if err != nil {
				t.Fatalf("unable to create testcase request")
			}
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			testutil.AssertEq(t, "status code", tc.status, res.Code)
		})
	}
}

func TestRoutingAccountGetBalance(t *testing.T) {
	tt := []struct {
		name      string
		service   func(t *testing.T) service.Account
		status    int
		path      string
		assertRes func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:   "get '/{id}/balance' successfully",
			status: http.StatusOK,
			path:   "/1/balance",
			service: func(t *testing.T) service.Account {
				return &testutil.AccountServMock{
					ExpectGetBalance: func(c context.Context, i int64) (float64, error) {
						testutil.AssertEq(t, "id", int64(1), i)
						return 50, nil
					},
				}
			},
			assertRes: func(t *testing.T, rec *httptest.ResponseRecorder) {
				testutil.AssertEq(t, "response balance", "50", string(bytes.TrimSpace(rec.Body.Bytes())))
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := chi.NewRouter()
			s := tc.service(t)
			r.Route("/", routing.Accounts(&s))

			req, err := http.NewRequest(http.MethodGet, tc.path, nil)
			if err != nil {
				t.Fatalf("unable to create testcase request")
			}
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			testutil.AssertEq(t, "status code", tc.status, res.Code)
			tc.assertRes(t, res)
		})
	}
}

func TestRoutingAccountCreate(t *testing.T) {
	tt := []struct {
		name      string
		service   func() service.Account
		status    int
		path      string
		assertRes func(*testing.T, *httptest.ResponseRecorder)
		reader    func() (io.Reader, error)
	}{
		{
			name:   "post '/' successfully",
			status: http.StatusCreated,
			path:   "/",
			service: func() service.Account {
				return &testutil.AccountServMock{
					ExpectCreate: func(c context.Context, ac dto.AccountCreation) (dto.AccountView, error) {
						return *testutil.NewAccountView(4, "Garcia", "33333333333", 25, time.Now()), nil
					},
				}
			},
			reader: func() (io.Reader, error) {
				if body, err := json.Marshal(testutil.NewAccountCreation("Garcia", "33333333333", "", 25)); err == nil {
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
			r.Route("/", routing.Accounts(&s))
			buffer, err := tc.reader()
			if err != nil {
				t.Fatalf("unable to create request body, %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, tc.path, buffer)
			if err != nil {
				t.Fatalf("unable to create testcase request")
			}
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			testutil.AssertEq(t, "status code", tc.status, res.Code)
		})
	}
}
