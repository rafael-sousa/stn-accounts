// Package entity groups types that models a database table
package entity

import (
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

// Constants related to Account fields
const (
	AccountNameSize   int = 255
	AccountCPFSize    int = 11
	AccountSecretSize int = 50
)

// Account models a financial account
type Account struct {
	ID        int64
	Name      string
	CPF       string
	Secret    string
	Balance   types.Currency
	CreatedAt time.Time
}
