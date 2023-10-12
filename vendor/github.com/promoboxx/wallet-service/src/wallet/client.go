package wallet

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/promoboxx/go-client/client"
	"github.com/promoboxx/go-discovery/src/discovery"
	"github.com/promoboxx/go-glitch/glitch"
)

// Query parameters for wallet endpoints
const (
	QueryParameterBrandID      = "brand_id"
	QueryParameterBusinessID   = "business_id"
	QueryParameterBusinessIDs  = "business_ids"
	QueryParameterRestrictions = "restrictions"
	QueryParameterExclusions   = "exclusions"
	QueryParameterFunded       = "funded"
)

//go:generate mockgen -source ./client.go -destination=./walletmock/client-mock.go -package=walletmock

// Client is a client that can interact with the service
type Client interface {
	GetWallets(ctx context.Context, authToken, walletType, ID string) (Wallets, glitch.DataError)
	GetLedgers(ctx context.Context, authToken, walletType, ID string, restrictions, exclusions []string, funded, excludeCampaigns bool) (Ledgers, glitch.DataError)
	MakeTransfer(ctx context.Context, authToken string, transferData TransferRequest) glitch.DataError
	ResetTransfer(ctx context.Context, authToken, idempotencyKey string) glitch.DataError
	MakeDeposit(ctx context.Context, authToken string, depositData DepositRequest) (DepositResponse, glitch.DataError)
	MakeWithdrawal(ctx context.Context, authToken string, withdrawalData WithdrawalRequest) (WithdrawalResponse, glitch.DataError)
	GetTransactionHistory(ctx context.Context, authToken, walletType, walletTypeID string, pageSize, offset *int32, sort *string) (ListTransactionResponse, glitch.DataError)
	GetTransactionHistoryFromKey(ctx context.Context, authToken, idempotencyKey string) (LedgerTransactionResponse, glitch.DataError)
	DeleteTransactionHistory(ctx context.Context, authToken, idempotencyKey string) glitch.DataError
}

func (sc *serviceClient) GetWallets(ctx context.Context, authToken, walletType, ID string) (Wallets, glitch.DataError) {
	var response Wallets
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	query := url.Values{}

	switch walletType {
	case TypeBrandWallet:
		query.Set(QueryParameterBrandID, ID)
	case TypeBusinessWallet:
		query.Set(QueryParameterBusinessID, ID)
	}

	err := sc.c.Do(
		ctx, http.MethodGet,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/wallets"),
		query, headers, nil, &response,
	)

	return response, err
}

func (sc *serviceClient) GetLedgers(ctx context.Context, authToken, walletType, ID string, restrictions, exclusions []string, funded, excludeCampaigns bool) (Ledgers, glitch.DataError) {
	var response Ledgers
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	query := url.Values{}

	switch walletType {
	case TypeBrandWallet:
		query.Set(QueryParameterBrandID, ID)
	case TypeBusinessWallet:
		query.Set(QueryParameterBusinessID, ID)
	}

	if funded {
		query.Set(QueryParameterFunded, "true")
	}

	if excludeCampaigns {
		exclusions = append(exclusions, "campaign:*")
	}

	if len(restrictions) > 0 {
		query.Set(QueryParameterRestrictions, strings.Join(restrictions, ","))
	}

	if len(exclusions) > 0 {
		query.Set(QueryParameterExclusions, strings.Join(exclusions, ","))
	}

	err := sc.c.Do(
		ctx, http.MethodGet,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/ledgers"),
		query, headers, nil, &response,
	)

	return response, err
}

func (sc *serviceClient) MakeTransfer(ctx context.Context, authToken string, transferData TransferRequest) glitch.DataError {
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	reader, err := client.ObjectToJSONReader(&transferData)
	if err != nil {
		return err
	}

	err = sc.c.Do(
		ctx, http.MethodPost,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/transfers"),
		nil, headers, reader, nil,
	)

	return err
}

func (sc *serviceClient) ResetTransfer(ctx context.Context, authToken, idempotencyKey string) glitch.DataError {
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	err := sc.c.Do(
		ctx, http.MethodPut,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/transfers/"+idempotencyKey+"/reset"),
		nil, headers, nil, nil,
	)

	return err
}

func (sc *serviceClient) MakeDeposit(ctx context.Context, authToken string, depositData DepositRequest) (DepositResponse, glitch.DataError) {
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	reader, err := client.ObjectToJSONReader(&depositData)
	if err != nil {
		return DepositResponse{}, err
	}

	var response DepositResponse

	err = sc.c.Do(
		ctx, http.MethodPost,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/deposits"),
		nil, headers, reader, &response,
	)

	return response, err
}

func (sc *serviceClient) MakeWithdrawal(ctx context.Context, authToken string, withdrawalData WithdrawalRequest) (WithdrawalResponse, glitch.DataError) {
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	reader, err := client.ObjectToJSONReader(&withdrawalData)
	if err != nil {
		return WithdrawalResponse{}, err
	}

	var response WithdrawalResponse

	err = sc.c.Do(
		ctx, http.MethodPost,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/withdrawals"),
		nil, headers, reader, &response,
	)

	return response, err
}

func (sc *serviceClient) GetTransactionHistory(ctx context.Context, authToken, walletType, walletTypeID string, pageSize, offset *int32, sort *string) (ListTransactionResponse, glitch.DataError) {
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	query := url.Values{}

	switch walletType {
	case TypeBrandWallet:
		query.Set(QueryParameterBrandID, walletTypeID)
	case TypeBusinessWallet:
		query.Set(QueryParameterBusinessID, walletTypeID)
	}

	if offset != nil {
		query.Set("offset", strconv.FormatInt(int64(*offset), 10))
	}

	if pageSize != nil {
		query.Set("page_size", strconv.FormatInt(int64(*pageSize), 10))
	}

	if sort != nil {
		query.Set("sort", *sort)
	}

	var response ListTransactionResponse

	err := sc.c.Do(
		ctx, http.MethodGet,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/transaction_history"),
		query, headers, nil, &response,
	)

	return response, err
}

func (sc *serviceClient) GetTransactionHistoryFromKey(ctx context.Context, authToken, idempotencyKey string) (LedgerTransactionResponse, glitch.DataError) {
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	var response LedgerTransactionResponse

	err := sc.c.Do(
		ctx, http.MethodGet,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/transaction_history/"+idempotencyKey),
		nil, headers, nil, &response,
	)

	return response, err
}

func (sc *serviceClient) DeleteTransactionHistory(ctx context.Context, authToken, idempotencyKey string) glitch.DataError {
	headers := http.Header{}
	headers.Add("Authorization", authToken)

	err := sc.c.Do(
		ctx, http.MethodDelete,
		client.PrefixRoute(sc.serviceName, sc.pathPrefix, sc.appendServiceNameToRoute, "v1/transaction_history/"+idempotencyKey),
		nil, headers, nil, nil,
	)

	return err
}

type serviceClient struct {
	c                        client.BaseClient
	pathPrefix               string
	appendServiceNameToRoute bool
	serviceName              string
}

// NewClient will create a new Client
func NewClient(finder discovery.Finder, useTLS bool, pathPrefix string, serviceName string, appendServiceNameToRoute bool, tlsConfig *tls.Config) Client {
	return &serviceClient{
		c:                        client.NewBaseClient(finder.FindService, serviceName, useTLS, 30*time.Second, tlsConfig),
		serviceName:              serviceName,
		appendServiceNameToRoute: appendServiceNameToRoute,
		pathPrefix:               pathPrefix,
	}
}
