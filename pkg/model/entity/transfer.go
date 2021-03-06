package entity

import (
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/types"
)

// Transfer registers a balance exchange between different accounts
type Transfer struct {
	ID          int64
	Origin      int64
	Destination int64
	Amount      types.Currency
	CreatedAt   time.Time
}
