package actorhub

import (
	"bytes"
	"context"
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
	// DefaultBaseURL is the default ActorHub API base URL.
	DefaultBaseURL = "https://api.actorhub.ai"

	// DefaultTimeout is the default request timeout.
	DefaultTimeout = 30 * time.Second

	// DefaultMaxRetries is the default number of retry attempts.
	DefaultMaxRetries = 3

	// Version is the SDK version.
	Version = "0.1.0"
)

// Client is the ActorHub API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	maxRetries int
}

// ClientOption is a function that configures the client.
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithTimeout sets a custom timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *Client) {
		c.maxRetries = maxRetries
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new ActorHub API client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		maxRetries: DefaultMaxRetries,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// doRequest performs an HTTP request with retry logic.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		err := c.doRequestOnce(ctx, method, path, body, result)
		if err == nil {
			return nil
		}

		lastErr = err

		// Only retry on rate limit or server errors
		switch err.(type) {
		case *RateLimitError, *ServerError:
			waitTime := time.Duration(1<<attempt) * time.Second
			if waitTime > 10*time.Second {
				waitTime = 10 * time.Second
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(waitTime):
				continue
			}
		default:
			return err
		}
	}

	return lastErr
}

// doRequestOnce performs a single HTTP request.
func (c *Client) doRequestOnce(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	reqURL := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "actorhub-go/"+Version)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, result)
}

// handleResponse processes the HTTP response.
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	requestID := resp.Header.Get("X-Request-ID")

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		var errResp map[string]interface{}
		json.Unmarshal(respBody, &errResp)
		message := "Invalid or missing API key"
		if detail, ok := errResp["detail"].(string); ok {
			message = detail
		}
		return NewAuthenticationError(message, requestID)
	}

	if resp.StatusCode == http.StatusNotFound {
		var errResp map[string]interface{}
		json.Unmarshal(respBody, &errResp)
		message := "Resource not found"
		if detail, ok := errResp["detail"].(string); ok {
			message = detail
		}
		return NewNotFoundError(message, requestID)
	}

	if resp.StatusCode == http.StatusUnprocessableEntity {
		var errResp map[string]interface{}
		json.Unmarshal(respBody, &errResp)
		message := "Validation error"
		if detail, ok := errResp["detail"].(string); ok {
			message = detail
		}
		errors, _ := errResp["errors"].(map[string]interface{})
		return NewValidationError(message, errors, requestID)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		var errResp map[string]interface{}
		json.Unmarshal(respBody, &errResp)
		message := "Rate limit exceeded"
		if detail, ok := errResp["detail"].(string); ok {
			message = detail
		}
		retryAfter := 0
		if ra := resp.Header.Get("Retry-After"); ra != "" {
			retryAfter, _ = strconv.Atoi(ra)
		}
		return NewRateLimitError(message, retryAfter, requestID)
	}

	if resp.StatusCode >= 500 {
		var errResp map[string]interface{}
		json.Unmarshal(respBody, &errResp)
		message := fmt.Sprintf("Server error: %d", resp.StatusCode)
		if detail, ok := errResp["detail"].(string); ok {
			message = detail
		}
		return NewServerError(message, resp.StatusCode, requestID)
	}

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.Unmarshal(respBody, &errResp)
		message := fmt.Sprintf("API error: %d", resp.StatusCode)
		if detail, ok := errResp["detail"].(string); ok {
			message = detail
		}
		return &ActorHubError{
			Message:      message,
			StatusCode:   resp.StatusCode,
			ResponseData: errResp,
			RequestID:    requestID,
		}
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// Verify checks if an image contains protected identities.
func (c *Client) Verify(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error) {
	if req.ImageURL == "" && req.ImageBase64 == "" {
		return nil, NewValidationError("Must provide image_url or image_base64", nil, "")
	}

	var result VerifyResponse
	err := c.doRequest(ctx, http.MethodPost, "/api/v1/identity/verify", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIdentity retrieves identity details by ID.
func (c *Client) GetIdentity(ctx context.Context, identityID string) (*IdentityResponse, error) {
	var result IdentityResponse
	err := c.doRequest(ctx, http.MethodGet, "/api/v1/identity/"+identityID, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CheckConsent checks consent status for face before AI generation.
func (c *Client) CheckConsent(ctx context.Context, req *ConsentCheckRequest) (*ConsentCheckResponse, error) {
	if req.ImageURL == "" && req.ImageBase64 == "" && len(req.FaceEmbedding) == 0 {
		return nil, NewValidationError("Must provide image_url, image_base64, or face_embedding", nil, "")
	}

	var result ConsentCheckResponse
	err := c.doRequest(ctx, http.MethodPost, "/api/v1/consent/check", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// ListMarketplace searches marketplace listings.
func (c *Client) ListMarketplace(ctx context.Context, req *MarketplaceListRequest) ([]MarketplaceListingResponse, error) {
	params := url.Values{}

	if req != nil {
		if req.Query != "" {
			params.Set("query", req.Query)
		}
		if req.Category != "" {
			params.Set("category", req.Category)
		}
		if len(req.Tags) > 0 {
			params.Set("tags", strings.Join(req.Tags, ","))
		}
		if req.Featured != nil {
			params.Set("featured", strconv.FormatBool(*req.Featured))
		}
		if req.MinPrice != nil {
			params.Set("min_price", strconv.FormatFloat(*req.MinPrice, 'f', -1, 64))
		}
		if req.MaxPrice != nil {
			params.Set("max_price", strconv.FormatFloat(*req.MaxPrice, 'f', -1, 64))
		}
		if req.SortBy != "" {
			params.Set("sort_by", req.SortBy)
		}
		if req.Page > 0 {
			params.Set("page", strconv.Itoa(req.Page))
		}
		if req.Limit > 0 {
			params.Set("limit", strconv.Itoa(req.Limit))
		}
	}

	path := "/api/v1/marketplace/listings"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result []MarketplaceListingResponse
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetMyLicenses retrieves licenses purchased by the current user.
func (c *Client) GetMyLicenses(ctx context.Context, status string, page, limit int) ([]LicenseResponse, error) {
	params := url.Values{}
	if status != "" {
		params.Set("status", status)
	}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}

	path := "/api/v1/marketplace/licenses/mine"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	var result []LicenseResponse
	err := c.doRequest(ctx, http.MethodGet, path, nil, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// PurchaseLicense purchases a license for an identity.
func (c *Client) PurchaseLicense(ctx context.Context, req *PurchaseLicenseRequest) (*PurchaseResponse, error) {
	if req.DurationDays == 0 {
		req.DurationDays = 30
	}

	var result PurchaseResponse
	err := c.doRequest(ctx, http.MethodPost, "/api/v1/marketplace/license/purchase", req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetActorPack retrieves Actor Pack status and details.
func (c *Client) GetActorPack(ctx context.Context, packID string) (*ActorPackResponse, error) {
	var result ActorPackResponse
	err := c.doRequest(ctx, http.MethodGet, "/api/v1/actor-packs/status/"+packID, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
