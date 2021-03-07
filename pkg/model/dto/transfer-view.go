package dto

import (
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
)

// TransferView exposes the displayable entity.Transfer values
type TransferView struct {
	ID          int64     `json:"id"`
	Destination int64     `json:"account_destination_id"`
	Amount      float64   `json:"amount"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewTransferView creates a view from the entity.Transfer stored at e
func NewTransferView(e *entity.Transfer) *TransferView {
	return &TransferView{
		ID:          e.ID,
		Destination: e.Destination,
		Amount:      e.Amount.Float64(),
		CreatedAt:   e.CreatedAt,
	}
}
