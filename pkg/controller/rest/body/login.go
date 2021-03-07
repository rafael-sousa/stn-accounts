package body

// LoginRequest holds the required fields for the login operation
type LoginRequest struct {
	CPF    string `json:"cpf" validation:"required" minLength:"11" maxLength:"11"`
	Secret string `json:"secret" validation:"required" minLength:"1" maxLength:"50"`
}

// LoginResponse maintains the response body of a successfull login
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}
