package wallet

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/promoboxx/go-service/uuid"
)

// Deposit source types
const (
	SourceContract = "contract"
	SourceStripe   = "stripe"
)

// BrandDeposit represents the data for a deposit
type BrandDeposit struct {
	ID            string    `json:"id"`
	WalletID      string    `json:"wallet_id"`
	Source        string    `json:"source"`
	SourceID      string    `json:"source_id"`
	AmountCents   int64     `json:"amount_cents,string"`
	TransactionID string    `json:"transaction_id,omitempty"`
	CreatedAt     time.Time `json:"-"`
	CreatedByID   string    `json:"-"`
	UpdatedAt     time.Time `json:"-"`
	UpdatedByID   string    `json:"-"`
}

// BrandDepositRequest represents the data sent to the API to be used in a deposit
type BrandDepositRequest struct {
	Source      string `json:"source"`
	SourceID    string `json:"source_id"`
	AmountCents int64  `json:"amount_cents,string"`
}

// Valid ensures the data represents a valid request
func (bdr BrandDepositRequest) Valid() (bool, RequestError) {
	if bdr.Source == "" {
		return false, RequestError{
			Detail:     "you must provide a source",
			Code:       ErrorInvalidSource,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(`brand deposit request requires a source`),
		}
	}

	if bdr.SourceID == "" {
		return false, RequestError{
			Detail:     "you must provide a source ID",
			Code:       ErrorInvalidSourceID,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(`brand deposit request requires a source ID`),
		}
	}

	if bdr.AmountCents < 1 {
		return false, RequestError{
			Detail:     "you must provide a non-zero, positive value for amount cents",
			Code:       ErrorInvalidAmountCents,
			Status:     http.StatusBadRequest,
			InnerError: fmt.Errorf(`amount cents "%d" for brand deposit request is invalid`, bdr.AmountCents),
		}
	}

	switch bdr.Source {
	case SourceContract:
		if !uuid.IsValid(bdr.SourceID) {
			return false, RequestError{
				Detail:     "invalid source ID",
				Code:       ErrorInvalidSourceID,
				Status:     http.StatusBadRequest,
				InnerError: fmt.Errorf(`source ID "%s" for brand deposit request is invalid`, bdr.SourceID),
			}
		}
	default:
		return false, RequestError{
			Detail:     "invalid source",
			Code:       ErrorInvalidSource,
			Status:     http.StatusBadRequest,
			InnerError: fmt.Errorf(`source "%s" for brand deposit request is invalid`, bdr.Source),
		}
	}

	return true, RequestError{}
}

// BrandDepositResponse represents the data sent back when a deposit request is made successfully
type BrandDepositResponse struct {
	DepositID     string `json:"deposit_id"`
	TransactionID string `json:"transaction_id"`
}
