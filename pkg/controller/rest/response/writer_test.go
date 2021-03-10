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
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
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
				testutil.AssertEq(t, "response body", 0, len(body))
			},
		},
		{
			name:       "write success response with location header",
			statusCode: http.StatusCreated,
			method:     http.MethodPost,
			resourceID: 1,
			assertResponse: func(t *testing.T, r *http.Response) {
				testutil.AssertEq(t, "location header", "/foo/1", r.Header.Get("Location"))
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
				testutil.AssertEq(t, "response body", `"content"`, string(bytes.TrimSpace(b)))
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

			testutil.AssertEq(t, "status code", tc.statusCode, res.Code)
			if res.Code == tc.statusCode {
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

			testutil.AssertEq(t, "status code", tc.statusCode, res.Code)
			if res.Code == tc.statusCode {
				b := body.JSONError{}
				if err = json.NewDecoder(res.Body).Decode(&b); err != nil {
					t.Errorf("unable to parse response body")
				}
			}
		})
	}
}
