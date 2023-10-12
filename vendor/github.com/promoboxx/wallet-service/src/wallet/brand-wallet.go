package wallet

import (
	"time"
)

// BrandWallet represents the data for a wallet
type BrandWallet struct {
	ID          string
	BrandID     string
	CreatedAt   time.Time
	CreatedByID string
	UpdatedAt   time.Time
	UpdatedByID string
}

type GetBrandRetailerLedgersRequest struct {
	BusinessIDs []string `json:"business_ids"`
}

// PaymentMethod represents the data for a brand wallet payment method
type PaymentMethod struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	BalanceCents int64  `json:"balance_cents,string"`
}

// BrandWalletResponse represents the data for a get brand wallet request
type BrandWalletResponse struct {
	ID             string          `json:"id"` // brand ID
	PaymentMethods []PaymentMethod `json:"payment_methods"`
}
