// Package dto holds types meant to carry values for an specific task, or to limit the amount of info exposed to outer application layers
package dto

// AccountCreation holds the values required for a entity.Account creation
type AccountCreation struct {
	Name    string  `json:"name" minLength:"1" maxLength:"255" example:"José da Silva" validate:"required"`
	CPF     string  `json:"cpf" minLength:"11" maxLength:"11" example:"11881200000"`
	Secret  string  `json:"secret" minLength:"1" maxLength:"50" example:"super_secret"`
	Balance float64 `json:"balance" minimum:"0"`
}
