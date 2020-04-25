package database

import "context"

type DBInstrumentTimer interface {
	End() error
}

// StartDBTimer starts a special DB timer and returns something that must be End()-ed
func StartDBTimer(ctx context.Context, collection, operation, query string) DBInstrumentTimer {
	return startNullDBTimer(ctx, collection, operation, query)
}

type nullDBTimer struct{}

func (nullDBTimer) End() error {
	return nil
}

func startNullDBTimer(ctx context.Context, collection, operation, query string) nullDBTimer {
	return nullDBTimer{}
}
