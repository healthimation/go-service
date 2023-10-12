package wallet

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/promoboxx/go-service/uuid"
)

// TransferData represents the data necessary to represent a source or destination for a transfer request
type TransferData struct {
	ID                  string      `json:"id"`
	Type                string      `json:"type"`
	TransferRestriction Restriction `json:"restriction"`
	ExpirationDate      *time.Time  `json:"expiration_date,omitempty"`
}

// Valid ensures the data for the TransferData struct is appropriate
func (td TransferData) Valid() (bool, RequestError) {
	if td.ID == "" || !uuid.IsValid(td.ID) {
		return false, RequestError{
			Detail:     "you must provide a valid ID as a UUID",
			Code:       ErrorInvalidTransferID,
			Status:     http.StatusBadRequest,
			InnerError: errors.New("wallet transfer request sent an invalid transfer ID [" + td.ID + "]"),
		}
	}

	switch td.Type {
	case TypeBrandWallet, TypeBusinessWallet:
	// no-op
	default:
		return false, RequestError{
			Detail:     "you must provide a valid type of either " + TypeBusinessWallet + " or " + TypeBrandWallet,
			Code:       ErrorInvalidTransferType,
			Status:     http.StatusBadRequest,
			InnerError: errors.New("wallet transfer request sent an invalid transfer type [" + td.Type + "]"),
		}
	}

	if valid, err := td.TransferRestriction.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid restriction",
			Code:       ErrorInvalidRestriction,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	return true, RequestError{}
}

// Transfer represents an individual transfer action to be taken
type Transfer struct {
	Source          TransferData    `json:"source"`
	Destination     TransferData    `json:"destination"`
	TransactionType string          `json:"transaction_type"`
	Metadata        json.RawMessage `json:"metadata"`
	AmountCents     int64           `json:"amount_cents,string"`
}

// Valid ensures the data for the Transfer struct is appropriate
func (t Transfer) Valid() (bool, RequestError) {
	if valid, err := t.Source.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid transfer source",
			Code:       ErrorInvalidTransferSource,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	if valid, err := t.Destination.Valid(); !valid {
		return false, RequestError{
			Detail:     "you must provide a valid transfer destination",
			Code:       ErrorInvalidTransferDestination,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(err.Error()),
		}
	}

	if t.AmountCents < 1 {
		return false, RequestError{
			Detail:     "you must provide a non-zero, positive value for transfer amount cents",
			Code:       ErrorInvalidAmountCents,
			Status:     http.StatusBadRequest,
			InnerError: fmt.Errorf(`amount cents "%d" for wallet transfer request is invalid`, t.AmountCents),
		}
	}

	if t.TransactionType == "" {
		return false, RequestError{
			Detail:     "transaction type must be a non-empty string",
			Code:       ErrorInvalidTransactionType,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(`empty transaction key provided for wallet transfer request`),
		}
	}

	return true, RequestError{}
}

// Value converts the transfer structure into a postgres compatible format
func (t Transfer) Value() (driver.Value, error) {
	ded := ""
	if t.Destination.ExpirationDate != nil {
		ded = t.Destination.ExpirationDate.Format(time.RFC3339)
	}

	var err error

	if t.Metadata == nil {
		// No metadata JSON was sent. Populate with an empty JSON object.
		t.Metadata = json.RawMessage(`{}`)
	} else {
		// Compact the JSON sent over. This strips out newlines, etc. that things like Postman beautify will send over.
		buffer := new(bytes.Buffer)
		err = json.Compact(buffer, t.Metadata)
		t.Metadata = buffer.Bytes()
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("(%s,%s,%s,%s,%s,%s,%s,%s,%s,%d,%s,%q)",
		t.Source.ID, t.Source.Type, t.Source.TransferRestriction.Key, t.Source.TransferRestriction.FilterValue,
		t.Destination.ID, t.Destination.Type, t.Destination.TransferRestriction.Key, t.Destination.TransferRestriction.FilterValue,
		ded, t.AmountCents, t.TransactionType, t.Metadata), nil
}

// TransferRequest represents the data sent to the API to be used to make a wallet transfer
// swagger:model TransferRequest
type TransferRequest struct {
	Transfers      []Transfer `json:"transfers"`
	IdempotencyKey string     `json:"idempotency_key"`
}

// Valid ensures the data for the TransferRequest struct is appropriate
func (tr TransferRequest) Valid() (bool, RequestError) {
	if !uuid.IsValid(tr.IdempotencyKey) {
		return false, RequestError{
			Detail:     "idempotency key must be a valid UUID",
			Code:       ErrorInvalidIdempotencyKey,
			Status:     http.StatusBadRequest,
			InnerError: fmt.Errorf(`idempotency key "%s" for wallet transfer request is invalid`, tr.IdempotencyKey),
		}
	}

	for _, t := range tr.Transfers {
		if valid, err := t.Valid(); !valid {
			return false, err
		}
	}

	return true, RequestError{}
}

// TransferRecord represents the data returned from the data layer for a particular transfer action
type TransferRecord struct {
	SourceType          string
	SourceLedgerID      string
	DestinationType     string
	DestinationLedgerID string
	AmountCents         int64
}
