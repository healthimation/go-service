package wallet

import "encoding/json"

type MarshalLedgers Ledgers

// Ledgers represents the data for a ledger collection
// swagger:model Ledgers
type Ledgers struct {
	Ledgers []Ledger `json:"ledgers"`
}

// MarshalJSON handles adding computed fields to the returned JSON structure.
func (ls Ledgers) MarshalJSON() ([]byte, error) {
	var totalAmtCents int64

	for _, l := range ls.Ledgers {
		totalAmtCents += l.Balance
	}

	return json.Marshal(struct {
		MarshalLedgers
		TotalAmountCents int64 `json:"total_amount_cents,string"`
	}{
		MarshalLedgers:   MarshalLedgers(ls),
		TotalAmountCents: totalAmtCents,
	})
}
