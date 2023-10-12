package wallet

import (
	"encoding/json"
	"github.com/promoboxx/go-service/service"
	"strings"
	"time"
)

// Wallet  types
const (
	TypeBrandWallet    = "brand"
	TypeBusinessWallet = "business"
)

// Wallet represents the data for a wallet
type Wallet struct {
	ID          string    `json:"id"`
	BrandID     string    `json:"brand_id"`
	BusinessID  string    `json:"business_id"`
	CreatedAt   time.Time `json:"-"`
	CreatedByID string    `json:"-"`
	Balance     int64     `json:"balance,string"`
	Ledgers     []Ledger  `json:"ledgers"`
}

// PagingDetails contains the total amount of records for a given query along with specific pagination details provided.
type PagingDetails struct {
	TotalResults int32   `json:"total_results"`
	PageSize     *int32  `json:"page_size,omitempty"`
	Offset       *int32  `json:"offset,omitempty"`
	Sort         *string `json:"sort,omitempty"`
}

// ParsePagingDetails takes the paging parameters and returns the details around them.
func ParsePagingDetails(pagingParams service.PagingParams) PagingDetails {
	pd := PagingDetails{}

	if pagingParams.Offset != nil {
		pd.Offset = pagingParams.Offset
	}

	if pagingParams.PageSize != nil {
		pd.PageSize = pagingParams.PageSize
	}

	if len(pagingParams.SortFields) > 0 {
		var s []string
		for _, sf := range pagingParams.SortFields {
			s = append(s, sf.Field+":"+sf.Direction)
		}
		sort := strings.Join(s, ",")
		pd.Sort = &sort
	}

	return pd
}

// TransactionHistory is the data for an individual record in a wallet's transaction history.
type TransactionHistory struct {
	ID                    string          `json:"id"`
	WalletID              string          `json:"wallet_id"`
	CreatedAt             time.Time       `json:"created_at"`
	CreatedByID           string          `json:"created_by_id"`
	SourceWalletType      string          `json:"source_wallet_type"`
	SourceWalletID        *string         `json:"source_wallet_id"`
	DestinationWalletType string          `json:"destination_wallet_type"`
	DestinationWalletID   *string         `json:"destination_wallet_id"`
	Amount                int64           `json:"amount"`
	IdempotencyKey        string          `json:"idempotency_key"`
	TransactionType       string          `json:"transaction_type"`
	Metadata              json.RawMessage `json:"metadata"`
	WalletRunningBalance  int64           `json:"wallet_running_balance"`
}

// ListTransactionResponse represents the response sent back for a "get transaction history" API call.
// swagger:model ListTransactionResponse
type ListTransactionResponse struct {
	Transactions []TransactionHistory `json:"transactions"`
	Paging       PagingDetails        `json:"paging"`
}
