package wallet

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/promoboxx/go-service/uuid"
)

// WithdrawalSource represents the data necessary for the source of a withdrawal request
type WithdrawalSource struct {
	ID                    string      `json:"id"`
	Type                  string      `json:"type"`
	WithdrawalRestriction Restriction `json:"restriction"`
}

// Valid ensures the data for the WithdrawalSource struct is appropriate
func (ws WithdrawalSource) Valid() (bool, RequestError) {
	if ws.ID == "" || !uuid.IsValid(ws.ID) {
		return false, RequestError{
			Detail:     "you must provide a valid ID as a UUID",
			Code:       ErrorInvalidWithdrawalID,
			Status:     http.StatusBadRequest,
			InnerError: errors.New("wallet withdrawal request sent an invalid withdrawal ID [" + ws.ID + "]"),
		}
	}

	switch ws.Type {
	case TypeBrandWallet, TypeBusinessWallet:
	// no-op
	default:
		return false, RequestError{
			Detail:     "you must provide a valid type of either " + TypeBusinessWallet + " or " + TypeBrandWallet,
			Code:       ErrorInvalidWithdrawalType,
			Status:     http.StatusBadRequest,
			InnerError: errors.New("wallet withdrawal request sent an invalid withdrawal type [" + ws.Type + "]"),
		}
	}

	if valid, err := ws.WithdrawalRestriction.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid restriction",
			Code:       ErrorInvalidRestriction,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	return true, RequestError{}
}

// WithdrawalDestination represents the data necessary for the destination of a withdrawal request
type WithdrawalDestination struct {
	WithdrawalRestriction Restriction `json:"restriction"`
}

// Valid ensures the data for the WithdrawalDestination struct is appropriate
func (wd WithdrawalDestination) Valid() (bool, RequestError) {
	if valid, err := wd.WithdrawalRestriction.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid restriction",
			Code:       ErrorInvalidRestriction,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	return true, RequestError{}
}

// Withdrawal represents an individual withdrawal action to be taken
type Withdrawal struct {
	Source          WithdrawalSource      `json:"source"`
	Destination     WithdrawalDestination `json:"destination"`
	TransactionType string                `json:"transaction_type"`
	Metadata        json.RawMessage       `json:"metadata"`
	AmountCents     int64                 `json:"amount_cents,string"`
}

// Valid ensures the data for the WithdrawalRequest struct is appropriate
func (w Withdrawal) Valid() (bool, RequestError) {
	if valid, err := w.Source.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid withdrawal source",
			Code:       ErrorInvalidWithdrawalSource,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	if valid, err := w.Destination.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid withdrawal destination",
			Code:       ErrorInvalidWithdrawalDestination,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	if w.AmountCents < 1 {
		return false, RequestError{
			Detail:     "you must provide a non-zero, positive value for withdrawal amount cents",
			Code:       ErrorInvalidAmountCents,
			Status:     http.StatusBadRequest,
			InnerError: fmt.Errorf(`amount cents "%d" for wallet withdrawal request is invalid`, w.AmountCents),
		}
	}

	if w.TransactionType == "" {
		return false, RequestError{
			Detail:     "transaction type must be a non-empty string",
			Code:       ErrorInvalidTransactionType,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(`empty transaction key provided for wallet withdrawal request`),
		}
	}

	return true, RequestError{}
}

// Value converts the withdrawal structure into a postgres compatible format
func (w Withdrawal) Value() (driver.Value, error) {
	var err error

	if w.Metadata == nil {
		// No metadata JSON was sent. Populate with an empty JSON object.
		w.Metadata = json.RawMessage(`{}`)
	} else {
		// Compact the JSON sent over. This strips out newlines, etc. that things like Postman beautify will send over.
		buffer := new(bytes.Buffer)
		err = json.Compact(buffer, w.Metadata)
		w.Metadata = buffer.Bytes()
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s,%s,%s,%s,%s,%s,%d,%s,%q)",
		w.Source.ID, w.Source.Type, w.Source.WithdrawalRestriction.Key, w.Source.WithdrawalRestriction.FilterValue,
		w.Destination.WithdrawalRestriction.Key, w.Destination.WithdrawalRestriction.FilterValue,
		w.AmountCents, w.TransactionType, w.Metadata), nil
}

// WithdrawalRequest represents the data sent to the API to be used to make a wallet withdrawal
// swagger:model WithdrawalRequest
type WithdrawalRequest struct {
	Withdrawals    []Withdrawal `json:"withdrawals"`
	IdempotencyKey string       `json:"idempotency_key"`
}

// Valid ensures the data for the WithdrawalRequest struct is appropriate
func (wr WithdrawalRequest) Valid() (bool, RequestError) {
	if !uuid.IsValid(wr.IdempotencyKey) {
		return false, RequestError{
			Detail:     "idempotency key must be a valid UUID",
			Code:       ErrorInvalidIdempotencyKey,
			Status:     http.StatusBadRequest,
			InnerError: fmt.Errorf(`idempotency key "%s" for wallet withdrawal request is invalid`, wr.IdempotencyKey),
		}
	}

	for _, w := range wr.Withdrawals {
		if valid, err := w.Valid(); !valid {
			return false, err
		}
	}

	return true, RequestError{}
}

// WithdrawalResponse represents the successful API response for a withdrawal request
// swagger:model WithdrawalResponse
type WithdrawalResponse struct {
	WalletID      string `json:"wallet_id"`
	WalletBalance int64  `json:"wallet_balance,string"`
}

// WithdrawalRecord represents the data returned from the data layer for a particular withdrawal action
type WithdrawalRecord struct {
	SourceType          string
	SourceLedgerID      string
	DestinationType     string
	DestinationLedgerID string
	AmountCents         int64
}
