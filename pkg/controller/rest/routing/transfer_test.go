package routing_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/routing"
	"github.com/rafael-sousa/stn-accounts/pkg/model/dto"
	"github.com/rafael-sousa/stn-accounts/pkg/service"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestRoutingTransferFetch(t *testing.T) {
	token, _, _ := jwtHandler.Generate(1)
	tt := []struct {
		name    string
		service func() service.Transfer
		status  int
		path    string
		headers map[string]string
	}{
		{
			name:   "get '/' without auth header",
			status: http.StatusUnauthorized,
			path:   "/",
			service: func() service.Transfer {
				return &testutil.TransferServMock{}
			},
		},
		{
			name:   "get '/' with no results successfully",
			status: http.StatusOK,
			path:   "/",
			service: func() service.Transfer {
				return &testutil.TransferServMock{
					ExpectFetch: func(c context.Context, i int64) ([]dto.TransferView, error) {
						testutil.AssertEq(t, "id", int64(1), i)
						return []dto.TransferView{}, nil
					},
				}
			},
			headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
		},
		{
			name:   "get '/' with results successfully",
			status: http.StatusOK,
			path:   "/",
			service: func() service.Transfer {
				return &testutil.TransferServMock{
					ExpectFetch: func(c context.Context, i int64) ([]dto.TransferView, error) {
						testutil.AssertEq(t, "id", int64(1), i)
						return []dto.TransferView{
							*testutil.NewTransferView(1, 2, 5),
							*testutil.NewTransferView(2, 2, 10),
							*testutil.NewTransferView(3, 2, 20),
						}, nil
					},
				}
			},
			headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := chi.NewRouter()
			s := tc.service()
			r.Route("/", routing.Transfers(&s, jwtHandler))

			req, err := http.NewRequest(http.MethodGet, tc.path, nil)
			if err != nil {
				t.Fatalf("unable to create testcase request")
			}
			if tc.headers != nil {
				header := req.Header
				for k, v := range tc.headers {
					header.Add(k, v)
				}
			}
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			testutil.AssertEq(t, "status code", tc.status, res.Code)
		})
	}
}

func TestRoutingTransferCreate(t *testing.T) {
	token, _, _ := jwtHandler.Generate(1)
	tt := []struct {
		name      string
		service   func() service.Transfer
		status    int
		path      string
		headers   map[string]string
		assertRes func(*testing.T, *httptest.ResponseRecorder)
		reader    func() (io.Reader, error)
	}{
		{
			name:   "post '/' successfully",
			status: http.StatusCreated,
			path:   "/",
			service: func() service.Transfer {
				return &testutil.TransferServMock{
					ExpectCreate: func(c context.Context, i int64, d dto.TransferCreation) (dto.TransferView, error) {
						testutil.AssertEq(t, "id", int64(1), i)
						testutil.AssertEq(t, "destination", int64(2), d.Destination)
						testutil.AssertEq(t, "amount", float64(500), d.Amount)
						return *testutil.NewTransferView(1, d.Destination, d.Amount), nil
					},
				}
			},
			headers: map[string]string{
				"Authorization": "Bearer " + token,
			},
			reader: func() (io.Reader, error) {
				if body, err := json.Marshal(testutil.NewTransferCreation(2, 500)); err == nil {
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
			r.Route("/", routing.Transfers(&s, jwtHandler))

			buffer, err := tc.reader()
			if err != nil {
				t.Fatalf("unable to create request body")
			}
			req, err := http.NewRequest(http.MethodPost, tc.path, buffer)
			if err != nil {
				t.Fatalf("unable to create testcase request")
			}
			if tc.headers != nil {
				header := req.Header
				for k, v := range tc.headers {
					header.Add(k, v)
				}
			}
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			testutil.AssertEq(t, "status code", tc.status, res.Code)
		})
	}
}
