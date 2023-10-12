package wallet

import "fmt"

// service level errors
const (
	ErrorDataStore             = "ERROR_DATA_STORE"
	ErrorForbidden             = "FORBIDDEN"
	ErrorInvalidJSON           = "INVALID_JSON"
	ErrorInvalidQueryParameter = "INVALID_QUERY_PARAMETER"
	ErrorNotAuthorized         = "NOT_AUTHORIZED"
	ErrorPermissionDenied      = "PERMISSION_DENIED"
	ErrorService               = "ERROR_SERVICE"
)

// errors specific to wallet transfers
const (
	ErrorInvalidTransferID          = "INVALID_TRANSFER_ID"
	ErrorInvalidTransferType        = "INVALID_TRANSFER_TYPE"
	ErrorInvalidTransferSource      = "INVALID_TRANSFER_SOURCE"
	ErrorInvalidTransferDestination = "INVALID_TRANSFER_DESTINATION"
)

// errors specific to wallet deposits
const (
	ErrorInvalidDepositSource      = "INVALID_DEPOSIT_SOURCE"
	ErrorInvalidDepositDestination = "INVALID_DEPOSIT_DESTINATION"
	ErrorInvalidDepositID          = "INVALID_DEPOSIT_ID"
	ErrorInvalidDepositType        = "INVALID_DEPOSIT_TYPE"
)

// errors specific to wallet withdrawals
const (
	ErrorInvalidWithdrawalSource      = "INVALID_WITHDRAWAL_SOURCE"
	ErrorInvalidWithdrawalDestination = "INVALID_WITHDRAWAL_DESTINATION"
	ErrorInvalidWithdrawalID          = "INVALID_WITHDRAWAL_ID"
	ErrorInvalidWithdrawalType        = "INVALID_WITHDRAWAL_TYPE"
)

// errors common to all wallet transactions
const (
	ErrorInvalidAmountCents     = "INVALID_AMOUNT_CENTS"
	ErrorInvalidBrandID         = "INVALID_BRAND_ID"
	ErrorInvalidBusinessID      = "INVALID_BUSINESS_ID"
	ErrorInvalidSource          = "INVALID_SOURCE"
	ErrorInvalidSourceID        = "INVALID_SOURCE_ID"
	ErrorInvalidSourceType      = "INVALID_SOURCE_TYPE"
	ErrorInvalidDestinationType = "INVALID_DESTINATION_TYPE"
	ErrorInvalidWalletType      = "INVALID_WALLET_TYPE"
	ErrorTooManyFilters         = "TOO_MANY_FILTERS"
	ErrorNoFilters              = "NO_FILTERS"
	ErrorInvalidRestriction     = "INVALID_RESTRICTION"
	ErrorInvalidExclusion       = "INVALID_EXCLUSION"
	ErrorInsufficientFunds      = "INSUFFICIENT_FUNDS"
	ErrorSourceWalletNotFound   = "SOURCE_WALLET_NOT_FOUND"
	ErrorInvalidIdempotencyKey  = "INVALID_IDEMPOTENCY_KEY"
	ErrorInvalidTransactionType = "INVALID_TRANSACTION_TYPE"
)

// errors for session lock work
const (
	ErrorScanningTask = "ERROR_SCANNING_TASK" // ErrorScanningTask is the code for errors scanning tasks for go-session-lock
)

// RequestError wraps the details of an HTTP request error
type RequestError struct {
	Detail     string
	Code       string
	Status     int
	InnerError error
}

// Error translates the error details into a single string
func (r RequestError) Error() string {
	return fmt.Sprintf("Code: [%s] Message: [%s] Inner error: [%s]", r.Code, r.Detail, r.InnerError)
}
