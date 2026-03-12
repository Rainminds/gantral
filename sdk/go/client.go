package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
func (c *Client) CreateInstance(ctx context.Context, workflowID string, triggerContext map[string]interface{}, pol Policy) (*Instance, error) {
	// Reconstruct the request body expected by the API
	reqPayload := map[string]interface{}{
		"workflow_id":     workflowID,
		"trigger_context": triggerContext,
		"policy":          pol,
	}

	bodyBytes, err := json.Marshal(reqPayload)
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

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &instance, nil
}

// RecordDecision records a human decision.
func (c *Client) RecordDecision(ctx context.Context, instanceID string, decisionType DecisionType, actorID, justification string) (*Instance, error) {
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

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		// Note that RecordDecision in the HTTP handler returns {"status": "SIGNAL_SENT"}
		// but the signature expects *Instance. In a real SDK we might provide a GetInstance long poll.
		// For now we'll allow decode to fail gracefully or return an empty instance if it's just a status ack.
		return &instance, nil
	}

	return &instance, nil
}

// GetInstance retrieves the current state of an instance.
func (c *Client) GetInstance(ctx context.Context, instanceID string) (*Instance, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/instances/%s", c.baseURL, instanceID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var instance Instance
	if err := json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &instance, nil
}
