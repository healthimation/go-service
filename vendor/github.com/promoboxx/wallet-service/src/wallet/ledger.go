package wallet

import (
	"encoding/json"
	"time"
)

type MarshalLedger Ledger

// Ledger represents the data for a ledger
type Ledger struct {
	ID             string        `json:"id"`
	WalletID       string        `json:"wallet_id"`
	Balance        int64         `json:"balance,string"`
	CreatedAt      time.Time     `json:"-"`
	CreatedByID    string        `json:"-"`
	Restrictions   []Restriction `json:"restrictions,omitempty"`
	ExpirationDate time.Time     `json:"expiration_date"`
	RefundedAt     time.Time     `json:"refunded_at"`
}

// MarshalJSON handles adding computed fields to the returned JSON structure.
func (l Ledger) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		MarshalLedger
		Expires  bool `json:"expires"`
		Refunded bool `json:"refunded"`
	}{
		MarshalLedger: MarshalLedger(l),
		Expires:       !l.ExpirationDate.IsZero(),
		Refunded:      !l.RefundedAt.IsZero(),
	})
}

// LedgerTransactionHistory is the data for an individual record in a wallet's transaction history, including restri.
type LedgerTransactionHistory struct {
	ID                     string          `json:"audit_log_action_id"`
	CreatedAt              time.Time       `json:"created_at"`
	CreatedByID            string          `json:"created_by_id"`
	SourceWalletType       string          `json:"source_wallet_type"`
	SourceWalletID         *string         `json:"source_wallet_id"`
	SourceRestriction      Restriction     `json:"source_restriction,omitempty"`
	SourceExpriation       time.Time       `json:"source_expiration_date"`
	DestinationWalletType  string          `json:"destination_wallet_type"`
	DestinationWalletID    *string         `json:"destination_wallet_id"`
	DestinationRestriction Restriction     `json:"destination_restriction,omitempty"`
	DestinationExpriation  time.Time       `json:"destination_expiration_date"`
	Amount                 int64           `json:"amount"`
	IdempotencyKey         string          `json:"idempotency_key"`
	TransactionType        string          `json:"transaction_type"`
	Metadata               json.RawMessage `json:"metadata"`
}

// LedgerTransactionResponse represents the response sent back for a "get transaction history for idempotency key" API call.
// swagger:model LedgerTransactionResponse
type LedgerTransactionResponse struct {
	Transactions []LedgerTransactionHistory `json:"transactions"`
}
