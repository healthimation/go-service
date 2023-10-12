package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/promoboxx/go-service/uuid"
)

// DepositSource represents the data necessary for the source of a deposit request
type DepositSource struct {
	DepositRestriction Restriction `json:"restriction"`
}

// Valid ensures the data for the DepositSource struct is appropriate
func (ds DepositSource) Valid() (bool, RequestError) {
	if valid, err := ds.DepositRestriction.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid restriction",
			Code:       ErrorInvalidRestriction,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	return true, RequestError{}
}

// DepositDestination represents the data necessary for the destination of a deposit request
type DepositDestination struct {
	ID                 string      `json:"id"`
	Type               string      `json:"type"`
	DepositRestriction Restriction `json:"restriction"`
	ExpirationDate     *time.Time  `json:"expiration_date,omitempty"`
}

// Valid ensures the data for the DepositDestination struct is appropriate
func (dd DepositDestination) Valid() (bool, RequestError) {
	if dd.ID == "" || !uuid.IsValid(dd.ID) {
		return false, RequestError{
			Detail:     "you must provide a valid ID as a UUID",
			Code:       ErrorInvalidDepositID,
			Status:     http.StatusBadRequest,
			InnerError: errors.New("wallet deposit request sent an invalid deposit ID [" + dd.ID + "]"),
		}
	}

	switch dd.Type {
	case TypeBrandWallet, TypeBusinessWallet:
	// no-op
	default:
		return false, RequestError{
			Detail:     "you must provide a valid type of either " + TypeBusinessWallet + " or " + TypeBrandWallet,
			Code:       ErrorInvalidDepositType,
			Status:     http.StatusBadRequest,
			InnerError: errors.New("wallet deposit request sent an invalid deposit type [" + dd.Type + "]"),
		}
	}

	if valid, err := dd.DepositRestriction.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid restriction",
			Code:       ErrorInvalidRestriction,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	return true, RequestError{}
}

// DepositRequest represents the data sent to the API to be used to make a wallet deposit
// swagger:model DepositRequest
type DepositRequest struct {
	Source          DepositSource      `json:"source"`
	Destination     DepositDestination `json:"destination"`
	TransactionType string             `json:"transaction_type"`
	Metadata        json.RawMessage    `json:"metadata"`
	AmountCents     int64              `json:"amount_cents,string"`
	IdempotencyKey  string             `json:"idempotency_key"`
}

// Valid ensures the data for the DepositRequest struct is appropriate
func (dr DepositRequest) Valid() (bool, RequestError) {
	if valid, err := dr.Source.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid deposit source",
			Code:       ErrorInvalidDepositSource,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	if valid, err := dr.Destination.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid deposit destination",
			Code:       ErrorInvalidDepositDestination,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	if dr.AmountCents < 1 {
		return false, RequestError{
			Detail:     "you must provide a non-zero, positive value for deposit amount cents",
			Code:       ErrorInvalidAmountCents,
			Status:     http.StatusBadRequest,
			InnerError: fmt.Errorf(`amount cents "%d" for wallet deposit request is invalid`, dr.AmountCents),
		}
	}

	if !uuid.IsValid(dr.IdempotencyKey) {
		return false, RequestError{
			Detail:     "idempotency key must be a valid UUID",
			Code:       ErrorInvalidIdempotencyKey,
			Status:     http.StatusBadRequest,
			InnerError: fmt.Errorf(`idempotency key "%s" for wallet deposit request is invalid`, dr.IdempotencyKey),
		}
	}

	if dr.TransactionType == "" {
		return false, RequestError{
			Detail:     "transaction type must be a non-empty string",
			Code:       ErrorInvalidTransactionType,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(`empty transaction key provided for wallet deposit request`),
		}
	}

	return true, RequestError{}
}

// DepositResponse represents the successful API response for a deposit request
// swagger:model DepositResponse
type DepositResponse struct {
	WalletID      string `json:"wallet_id"`
	WalletBalance int64  `json:"wallet_balance,string"`
}
