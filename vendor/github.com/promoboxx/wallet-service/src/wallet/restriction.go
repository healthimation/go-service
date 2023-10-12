package wallet

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net/http"
)

// Restriction represents the data for a ledger restriction
type Restriction struct {
	Key         string `json:"key"`
	FilterValue string `json:"value"`
}

// Valid ensures the data represents a valid request
func (r Restriction) Valid() (bool, RequestError) {
	if r.Key == "" {
		return false, RequestError{
			Detail:     "you must provide a restriction key",
			Code:       ErrorInvalidRestriction,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(`restriction requires a key`),
		}
	}

	if r.FilterValue == "" {
		return false, RequestError{
			Detail:     "you must provide a restriction value",
			Code:       ErrorInvalidRestriction,
			Status:     http.StatusBadRequest,
			InnerError: errors.New(`restriction requires a value`),
		}
	}

	return true, RequestError{}
}

// Value converts the restriction filter into a postgres compatible format
func (r Restriction) Value() (driver.Value, error) {
	return fmt.Sprintf(`(%s,%s)`, r.Key, r.FilterValue), nil
}
