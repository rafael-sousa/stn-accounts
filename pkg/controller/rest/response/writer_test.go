package response_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/body"
	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/response"
	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

func TestWriteSuccess(t *testing.T) {
	tt := []struct {
		name           string
		responseBody   interface{}
		resourceID     interface{}
		assertResponse func(*testing.T, *http.Response)
		statusCode     int
		method         string
	}{
		{
			name:       "write success response without content",
			statusCode: http.StatusNoContent,
			method:     http.MethodGet,
			assertResponse: func(t *testing.T, r *http.Response) {
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Errorf("unable to read response body, %v", err)
					return
				}
				if len(body) != 0 {
					t.Errorf("expected response body size equal to '%d' but got '%d'", 0, len(body))
				}
			},
		},
		{
			name:       "write success response with location header",
			statusCode: http.StatusCreated,
			method:     http.MethodPost,
			resourceID: 1,
			assertResponse: func(t *testing.T, r *http.Response) {
				location := r.Header.Get("Location")
				if location != "/foo/1" {
					t.Errorf("expected location header equal to '%s' but got '%s'", "/foo/1", location)
				}
			},
			responseBody: "body",
		},
		{
			name:       "write success response with body content",
			statusCode: http.StatusOK,
			method:     http.MethodGet,
			assertResponse: func(t *testing.T, r *http.Response) {
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Errorf("unable to read response body, %v", err)
					return
				}
				if string(bytes.TrimSpace(b)) != `"content"` {
					t.Errorf("expected response body equal to '%s' but got '%s'", `"content"`, string(b))
				}
			},
			responseBody: "content",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(tc.method, "/foo", nil)
			if err != nil {
				t.Errorf("unabled to create http request, %v", err)
				return
			}
			res := httptest.NewRecorder()
			response.WriteSuccess(res, request, tc.responseBody, tc.resourceID)

			if res.Code != tc.statusCode {
				t.Errorf("expected status code equal to '%d' but got '%d'", tc.statusCode, res.Code)
			} else {
				tc.assertResponse(t, res.Result())
			}
		})
	}
}

func TestWriteErr(t *testing.T) {
	tt := []struct {
		name       string
		statusCode int
		err        error
	}{
		{
			name:       "write response with unknown error type",
			statusCode: http.StatusInternalServerError,
			err:        fmt.Errorf("unknown"),
		},
		{
			name:       "write response with NotFoundErr error type",
			statusCode: http.StatusNotFound,
			err:        types.NewErr(types.NotFoundErr, "NotFoundErr", nil),
		},
		{
			name:       "write response with EmptyResultErr error type",
			statusCode: http.StatusNotFound,
			err:        types.NewErr(types.EmptyResultErr, "EmptyResultErr", nil),
		},
		{
			name:       "write response with ValidationErr error type",
			statusCode: http.StatusBadRequest,
			err:        types.NewErr(types.ValidationErr, "ValidationErr", nil),
		},
		{
			name:       "write response with AuthenticationErr error type",
			statusCode: http.StatusUnauthorized,
			err:        types.NewErr(types.AuthenticationErr, "AuthenticationErr", nil),
		},
		{
			name:       "write response with ConflictErr error type",
			statusCode: http.StatusConflict,
			err:        types.NewErr(types.ConflictErr, "ConflictErr", nil),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "/foo", nil)
			if err != nil {
				t.Errorf("unabled to create http request, %v", err)
				return
			}
			res := httptest.NewRecorder()
			response.WriteErr(res, request, tc.err)

			if res.Code != tc.statusCode {
				t.Errorf("expected status code equal to '%d' but got '%d'", tc.statusCode, res.Code)
			} else {
				b := body.JSONError{}
				if err = json.NewDecoder(res.Body).Decode(&b); err != nil {
					t.Errorf("unable to parse response body")
				}
			}
		})
	}
}
