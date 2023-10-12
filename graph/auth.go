package graph

import (
	"context"
	"errors"

	"github.com/promoboxx/go-auth/src/auth"
	"github.com/promoboxx/go-service/alice/middleware"
	"github.com/promoboxx/wallet-service/src/wallet"
)

// This is a bit of a hack workaround to make the `middleware.GetClaimsFromContext` function return an error,
// to gracefully handle situations where no auth claims were injected by any middleware. This change is necessary
// because the GraphQL middleware does not prevent requests from reaching the resolvers if no token is present,
// so they could panic with an unhelpful error message if they tried to call the standard function
func EnsureAuthenticated(ctx context.Context) (c auth.Claim, e error) {
	defer func() {
		if err := recover(); err != nil {
			e = errors.New("unauthorized access")
		}
	}()

	return middleware.GetClaimsFromContext(ctx), nil
}

func IsAuthorized(claim auth.Claim, authType, ID string) bool {
	switch authType {
	case wallet.TypeBusinessWallet:
		return claim.HasPermission(auth.CheckBusinessID{BusinessID: ID}) || claim.IsInternal()
	case wallet.TypeBrandWallet:
		return claim.HasPermission(auth.CheckBrandUUID{BrandUUID: ID}) || claim.IsInternal()
	default:
		return claim.IsInternal()
	}
}
