package routing_test

import (
	"os"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/controller/rest/jwt"
	"github.com/rafael-sousa/stn-accounts/pkg/model/env"
)

var jwtHandler *jwt.Handler

func TestMain(m *testing.M) {
	jwtHandler = jwt.NewHandler(&env.RestConfig{
		Secret:          []byte("secret"),
		TokenExpTimeout: 30,
	})
	os.Exit(m.Run())
}
