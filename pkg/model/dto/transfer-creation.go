package dto

// TransferCreation holds the values required for a entity.Transfer creation
type TransferCreation struct {
	Destination int64   `json:"account_destination_id" validation:"required" minimum:"1"`
	Amount      float64 `json:"amount" validation:"required" minimum:"0.01"`
}
