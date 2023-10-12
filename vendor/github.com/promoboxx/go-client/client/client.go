package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/promoboxx/go-glitch/glitch"
	"github.com/promoboxx/go-service/alice/middleware"
)

// Error codes
const (
	ErrorCantFind          = "CANT_FIND_SERVICE"
	ErrorRequestCreation   = "CANT_CREATE_REQUEST"
	ErrorRequestError      = "ERROR_MAKING_REQUEST"
	ErrorDecodingError     = "ERROR_DECODING_ERROR"
	ErrorDecodingResponse  = "ERROR_DECODING_RESPONSE"
	ErrorMarshallingObject = "ERROR_MARSHALLING_OBJECT"
	ErrorURL               = "ERROR_URL"
)

// ServiceFinder can find a service's base URL
type ServiceFinder func(serviceName string) (string, error)

// BaseClient can do requests
//go:generate mockgen -destination=./clientmock/client-mock.go -package=clientmock github.com/promoboxx/go-client/client BaseClient
type BaseClient interface {
	// Do does the request and parses the body into the response provider if in the 2xx range, otherwise parses it into a glitch.DataError
	Do(ctx context.Context, method string, slug string, query url.Values, headers http.Header, body io.Reader, response interface{}) glitch.DataError

	// MakeRequest does the request and returns the status, body, and any error
	// This should be used only if the api doesn't return glitch.DataErrors
	MakeRequest(ctx context.Context, method string, slug string, query url.Values, headers http.Header, body io.Reader) (int, []byte, glitch.DataError)
}

type client struct {
	finder      ServiceFinder
	useTLS      bool
	serviceName string
	client      *http.Client
}

// NewBaseClient creates a new BaseClient
func NewBaseClient(finder ServiceFinder, serviceName string, useTLS bool, timeout time.Duration, tlsConfig *tls.Config) BaseClient {
	rt := http.DefaultTransport

	if useTLS && tlsConfig != nil {
		rt = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	c := &http.Client{
		Timeout:   timeout,
		Transport: rt,
	}

	return &client{finder: finder, serviceName: serviceName, useTLS: useTLS, client: c}
}

// Do does the request and parses the body into the response provider if in the 2xx range, otherwise parses it into a glitch.DataError
func (c *client) Do(ctx context.Context, method string, slug string, query url.Values, headers http.Header, body io.Reader, response interface{}) glitch.DataError {
	return c.do(ctx, method, slug, query, headers, body, response, nil)
}

// Do does the request and parses the body into the response provider if in the 2xx range, otherwise parses it into a glitch.DataError
// The name arg will be used to assign the name of the span that is in the client
func (c *client) DoWithName(ctx context.Context, method string, slug string, query url.Values, headers http.Header, body io.Reader, response interface{}, name string) glitch.DataError {
	return c.do(ctx, method, slug, query, headers, body, response, &name)
}

// MakeRequest does the request and returns the status, body, and any error
// This should be used only if the api doesn't return glitch.DataErrors
func (c *client) MakeRequest(ctx context.Context, method string, slug string, query url.Values, headers http.Header, body io.Reader) (int, []byte, glitch.DataError) {
	return c.makeRequest(ctx, method, slug, query, headers, body, nil)
}

// MakeRequest does the request and returns the status, body, and any error
// This should be used only if the api doesn't return glitch.DataErrors
// The name arg will be used to assign the name of the span that is in the client
func (c *client) MakeRequestWithName(ctx context.Context, method string, slug string, query url.Values, headers http.Header, body io.Reader, name string) (int, []byte, glitch.DataError) {
	return c.makeRequest(ctx, method, slug, query, headers, body, &name)
}

func (c *client) do(ctx context.Context, method string, slug string, query url.Values, headers http.Header, body io.Reader, response interface{}, name *string) glitch.DataError {
	var (
		status int
		ret    []byte
		err    glitch.DataError
	)

	if name != nil {
		status, ret, err = c.MakeRequestWithName(ctx, method, slug, query, headers, body, *name)
	} else {
		status, ret, err = c.MakeRequest(ctx, method, slug, query, headers, body)
	}

	if err != nil {
		return err
	}

	if status >= 400 || status < 200 {
		prob := glitch.HTTPProblem{}
		err := json.Unmarshal(ret, &prob)
		if err != nil {
			return glitch.NewDataError(err, ErrorRequestError, "Could not decode error response")
		}
		return glitch.FromHTTPProblem(prob, fmt.Sprintf("Error from %s to %s - %s", method, c.serviceName, slug))
	}

	if response != nil {
		err := json.Unmarshal(ret, response)
		if err != nil {
			return glitch.NewDataError(err, ErrorDecodingResponse, "Could not decode response")
		}
	}

	return nil
}

func (c *client) makeRequest(ctx context.Context, method string, slug string, query url.Values, headers http.Header, body io.Reader, name *string) (int, []byte, glitch.DataError) {
	rawURL, err := c.finder(c.serviceName)
	if err != nil {
		return 0, nil, glitch.NewDataError(err, ErrorCantFind, "Error finding service")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return 0, nil, glitch.NewDataError(err, ErrorURL, "Error parsing url from string")
	}

	u.Path = slug
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return 0, nil, glitch.NewDataError(err, ErrorRequestCreation, "Error creating request object")
	}

	if headers != nil {
		req.Header = headers
	} else {
		req.Header = http.Header{}
	}

	if ctx != nil {
		span := opentracing.SpanFromContext(ctx)

		// check if a name was set, if there was no name set
		// default to method and service name
		// span name cannot be more than 100 characters
		if name == nil {
			tmp := fmt.Sprintf("%s %s", method, c.serviceName)
			name = &tmp
		}

		// create the child span that will correlate to the
		// parent span if one exists
		var childSpan opentracing.Span
		if span != nil {
			childSpan = opentracing.StartSpan(*name, opentracing.ChildOf(span.Context()))
			defer childSpan.Finish()
			opentracing.GlobalTracer().Inject(childSpan.Context(), opentracing.HTTPHeaders, req.Header)
		} else {
			span = opentracing.StartSpan(*name)
			defer span.Finish()
		}

		span.SetTag("slug", slug)

		req = req.WithContext(ctx)

		// if we have a requestID in the context pass it along in the header
		requestID := middleware.GetRequestIDFromContext(ctx)
		if len(requestID) > 0 {
			req.Header.Set(middleware.HeaderRequestID, requestID)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, glitch.NewDataError(err, ErrorRequestError, "Could not make the request")
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, glitch.NewDataError(err, ErrorDecodingResponse, "Could not read response body")
	}

	return resp.StatusCode, ret, nil
}
