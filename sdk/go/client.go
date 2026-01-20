package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Rainminds/gantral/core/engine"
	"github.com/Rainminds/gantral/core/policy"
)

// Client is the Gantral SDK client.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Gantral client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// CreateInstance creates a new execution instance.
func (c *Client) CreateInstance(ctx context.Context, workflowID string, triggerContext map[string]interface{}, pol policy.Policy) (*engine.Instance, error) {
	// Reconstruct the request body expected by the API
	reqBody := triggerContext
	// Note: The API currently expects a generic map and parses "materiality" from it if extracting logic is in handlers.
	// The handler logic:
	// if mat, ok := body["materiality"].(string); ok && mat == "HIGH" { pol.Materiality = ... }
	// So we need to inject materiality into the map if we want to test that.
	// For Typed SDK, we might want to expose overrides.
	// For now, let's just pass the context map.

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/instances", c.baseURL), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var instance engine.Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &instance, nil
}

// RecordDecision records a human decision.
func (c *Client) RecordDecision(ctx context.Context, instanceID string, decisionType engine.DecisionType, actorID, justification string) (*engine.Instance, error) {
	reqBody := map[string]string{
		"type":          string(decisionType),
		"actor_id":      actorID,
		"justification": justification,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := fmt.Sprintf("%s/instances/%s/decisions", c.baseURL, instanceID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 201 Created from API
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var instance engine.Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &instance, nil
}
