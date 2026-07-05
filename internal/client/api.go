package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var DefaultClient *APIClient

type APIClient struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func New(basURL string) *APIClient {
	return &APIClient{
		BaseURL:    basURL,
		HTTPClient: http.DefaultClient,
	}
}

func (c *APIClient) SetToken(token string) {
	c.Token = token
}

func (c *APIClient) do(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

func (c *APIClient) Post(path string, body interface{}) ([]byte, int, error) {
	return c.do(http.MethodPost, path, body)
}

func (c *APIClient) Get(path string) ([]byte, int, error) {
	return c.do(http.MethodGet, path, nil)
}

func (c *APIClient) Delete(path string) ([]byte, int, error) {
	return c.do(http.MethodDelete, path, nil)
}
