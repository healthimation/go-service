package instrumenter

import (
	"context"
	"net/http"
)

// Client is a weak wrapper around the NR application 'start transaction' functionality, so we don't
// need to pull the NR lib directly into all the things.
type Client interface {
	// Start a top level transaction, omit w and r for a background transaction
	StartTransaction(ctx context.Context, name string, w http.ResponseWriter, r *http.Request) (context.Context, Transaction)
}

// Txn wraps the bits of the Newrelic Transaction interface that we use
type Transaction interface {
	End() error
	NoticeError(err error) error
}

// InstrumentTimer must be used within a Transaction, it maps to a Newrelic 'segment'
type InstrumentTimer interface {
	End() error
}

// ExternalCallTimer is a timer for service a to instrument outgoing calls within a Transaction
type ExternalCallTimer interface {
	End() // Stop the timer
}
