package pyembeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	EmbedPath = "/embed"
	SizePath  = "/size"
)

type EmbeddingRequest struct {
	Content string `json:"content"`
	Model   string `json:"model"`
	ApiKey  string `json:"api_key"`
}

type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

type SizeRequest struct {
}

type SizeResponse struct {
	Size int `json:"size"`
}

type Config struct {
	host       string
	httpClient *http.Client
}

func NewConfig(host string) *Config {
	return &Config{
		host: host,
	}
}

func (c *Config) WithHttpClient(httpClient *http.Client) *Config {
	c.httpClient = httpClient
	return c
}

type Client struct {
	client *http.Client
	url    *url.URL
	config *Config
}

func NewClient(ctx context.Context, config *Config) (*Client, error) {
	urlValue, err := url.Parse(config.host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	if config.httpClient == nil {
		config.httpClient = &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   30 * time.Second,
		}
	}

	context.AfterFunc(ctx, func() {
		config.httpClient.CloseIdleConnections()
	})

	return &Client{
		client: config.httpClient,
		config: config,
		url:    urlValue,
	}, nil
}

func (c *Client) Embed(content string) ([]float64, error) {
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(&EmbeddingRequest{Content: content})

	postReq := &http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme: c.url.Scheme,
			Host:   c.url.Host,
			Path:   EmbedPath,
		},
		Body: io.NopCloser(buffer),
	}

	postReq.Header = http.Header{
		"Content-Type": []string{"application/json"},
	}

	resp, err := c.client.Do(postReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var embedResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
		return nil, err
	}

	return embedResp.Embedding, nil
}

func (c *Client) Size() (int, error) {
	getReq := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: c.url.Scheme,
			Host:   c.url.Host,
		},
	}

	resp, err := c.client.Do(getReq)
	if err != nil {
		return 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var sizeResp SizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&sizeResp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return sizeResp.Size, nil
}
