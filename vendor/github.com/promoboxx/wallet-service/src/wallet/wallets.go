package wallet

// Wallets represents the data for a wallet collection
// swagger:model Wallets
type Wallets struct {
	Wallets []Wallet `json:"wallets"`
}
