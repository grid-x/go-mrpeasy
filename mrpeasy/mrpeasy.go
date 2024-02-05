package mrpeasy

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"

	headerContentRange = "content-range"
	headerRange        = "range"

	defaultTimeout    = 10 * time.Second
	defaultAPIBaseURL = "https://api.mrpeasy.com/rest/v1/"
)

type Client struct {
	client            *http.Client
	apiKey, apiSecret string
	apiBaseURL        *url.URL
}

type ClientOption func(*Client)

func WithHTTPClient(hcl *http.Client) ClientOption {
	return func(c *Client) {
		c.client = hcl
	}
}

func New(apiKey, apiSecret string, opts ...ClientOption) (*Client, error) {
	parsedURL, err := url.Parse(defaultAPIBaseURL)
	if err != nil {
		return nil, err
	}
	cl := &Client{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		apiBaseURL: parsedURL,
	}

	for _, opt := range opts {
		opt(cl)
	}

	return cl, nil
}

type Response struct {
	totalItems int
	lastItem   int
	*http.Response
}

func (r Response) HasNext() bool {
	return r.totalItems-1 > r.lastItem
}

func (r Response) Next() int {
	return r.lastItem
}

func resolveNextItem(r *http.Response) (*Response, error) {
	contentRange := r.Header.Get(headerContentRange)
	fullHeaderSplit := strings.Split(contentRange, "/")
	totalItems, err := strconv.Atoi(fullHeaderSplit[1])
	if err != nil {
		return nil, err
	}

	firstHeaderSplit := strings.Split(fullHeaderSplit[0], "-")
	lastItem, err := strconv.Atoi(firstHeaderSplit[1])
	if err != nil {
		return nil, err
	}

	return &Response{
		totalItems: totalItems,
		lastItem:   lastItem,
		Response:   r,
	}, nil
}

type RequestOption func(*http.Request)

// v of form "100" or "10-14"
func addRangeItemsHeader(r *http.Request, v string) {
	r.Header.Add(headerRange, fmt.Sprintf("items=%s", v))
}

// Request range of 100 items starting at from+1
func WithRangeFrom(from int) RequestOption {
	return func(r *http.Request) {
		addRangeItemsHeader(r, fmt.Sprintf("%d", from))
	}
}

// Request range of items starting at from+1 until to+1
func WithRangeFromTo(from, to int) RequestOption {
	return func(r *http.Request) {
		addRangeItemsHeader(r, fmt.Sprintf("%d-%d", from, to))
	}
}

func (c *Client) NewRequest(method string, url string, body interface{}, opts ...RequestOption) (*http.Request, error) {
	u, err := c.apiBaseURL.Parse(url)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Add basic authentication header
	authHeader := fmt.Sprintf("%s:%s", c.apiKey, c.apiSecret)
	b64AuthHeader := base64.StdEncoding.EncodeToString([]byte(authHeader))
	req.Header.Add("Authorization", "Basic "+b64AuthHeader)

	if body != nil {
		req.Header.Set(headerContentType, contentTypeJSON)
	}

	for _, o := range opts {
		o(req)
	}

	return req, nil
}

func isError(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	return fmt.Errorf("response: %s", r.Status)
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)

	// Perform the request
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	err = isError(resp)
	if err != nil {
		return &Response{Response: resp}, err
	}

	r, err := resolveNextItem(resp)
	if err != nil {
		return nil, err
	}

	if v != nil {
		// If v is io.Writer write response to it. Else json.Decode into v.
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, fmt.Errorf("could not copy respo.Body to writer: %+v", err)
			}
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return r, err
}

func (c *Client) ListCustomers(ctx context.Context) ([]*Customer, error) {
	var result []*Customer
	var opts []RequestOption
	for {
		req, err := c.NewRequest(http.MethodGet, "customers", nil, opts...)
		if err != nil {
			return nil, err
		}
		var tmp []*Customer
		resp, err := c.Do(ctx, req, &tmp)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp...)
		if !resp.HasNext() {
			break
		}

		opts = []RequestOption{WithRangeFrom(resp.Next())}
	}
	return result, nil
}

func (c *Client) ListCustomerOrders(ctx context.Context) ([]*CustomerOrder, error) {
	var result []*CustomerOrder
	var opts []RequestOption
	for {
		req, err := c.NewRequest(http.MethodGet, "customer-orders", nil, opts...)
		if err != nil {
			return nil, err
		}
		var tmp []*CustomerOrder
		resp, err := c.Do(ctx, req, &tmp)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp...)
		if !resp.HasNext() {
			break
		}

		opts = []RequestOption{WithRangeFrom(resp.Next())}
	}
	return result, nil
}

func (c *Client) ListShipments(ctx context.Context) ([]*Shipment, error) {
	var result []*Shipment
	var opts []RequestOption
	for {
		req, err := c.NewRequest(http.MethodGet, "shipments", nil, opts...)
		if err != nil {
			return nil, err
		}
		var tmp []*Shipment
		resp, err := c.Do(ctx, req, &tmp)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp...)
		if !resp.HasNext() {
			break
		}

		opts = []RequestOption{WithRangeFrom(resp.Next())}
	}
	return result, nil
}

func (c *Client) ListStockItems(ctx context.Context) ([]*StockItem, error) {
	var result []*StockItem
	var opts []RequestOption
	for {
		req, err := c.NewRequest(http.MethodGet, "items", nil, opts...)
		if err != nil {
			return nil, err
		}
		var tmp []*StockItem
		resp, err := c.Do(ctx, req, &tmp)
		if err != nil {
			return nil, err
		}
		result = append(result, tmp...)
		if !resp.HasNext() {
			break
		}

		opts = []RequestOption{WithRangeFrom(resp.Next())}
	}
	return result, nil
}
