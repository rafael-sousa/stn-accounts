package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

func TestRoutingNotfound(t *testing.T) {
	tt := []struct {
		name   string
		status int
		path   string
	}{
		{
			name:   "get '/foo' with no route",
			status: http.StatusNotFound,
			path:   "/foo",
		},
		{
			name:   "get '/bar' with no route",
			status: http.StatusNotFound,
			path:   "/bar",
		},
		{
			name:   "get '/' with no route",
			status: http.StatusNotFound,
			path:   "/",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := chi.NewRouter()
			req, err := http.NewRequest(http.MethodPost, tc.path, nil)
			if err != nil {
				t.Fatalf("unable to create testcase request")
			}
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			testutil.AssertEq(t, "status code", tc.status, res.Code)
		})
	}
}
