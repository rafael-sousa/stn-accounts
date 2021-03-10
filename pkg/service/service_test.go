package service_test

import (
	"os"
	"testing"

	"github.com/rafael-sousa/stn-accounts/pkg/repository"
	"github.com/rafael-sousa/stn-accounts/pkg/testutil"
)

var txr repository.Transactioner

func TestMain(m *testing.M) {
	txr = &testutil.TransactionerMock{}
	os.Exit(m.Run())
}
