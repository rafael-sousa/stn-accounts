package dto

import (
	"time"

	"github.com/rafael-sousa/stn-accounts/pkg/model/entity"
)

// AccountView maintains the displayable entity.Account values
type AccountView struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// NewAccountView creates a view from the entity.Account stored at e
func NewAccountView(e *entity.Account) *AccountView {
	return &AccountView{
		ID:        e.ID,
		Name:      e.Name,
		Balance:   e.Balance.Float64(),
		CPF:       e.CPF,
		CreatedAt: e.CreatedAt,
	}
}
