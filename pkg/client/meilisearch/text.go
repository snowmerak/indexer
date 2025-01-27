package meilisearch

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/meilisearch/meilisearch-go"

	"github.com/snowmerak/indexer/lib/index/text"
)

var _ text.Text = (*Client)(nil)

type Config struct {
	CollectionName           string
	Host                     string
	ApiKey                   string
	ContentEncoding          meilisearch.ContentEncoding
	EncodingCompressionLevel meilisearch.EncodingCompressionLevel
	HttpClient               *http.Client
	TlsConfig                *tls.Config
	RetryStatus              []int
	MaxRetries               int
}

func NewConfig(host string, collectionName string) *Config {
	return &Config{
		Host:           host,
		CollectionName: collectionName,
	}
}

func (c *Config) WithApiKey(apiKey string) *Config {
	c.ApiKey = apiKey
	return c
}

func (c *Config) WithContentEncoding(contentEncoding meilisearch.ContentEncoding, encodingCompressionLevel meilisearch.EncodingCompressionLevel) *Config {
	c.ContentEncoding = contentEncoding
	c.EncodingCompressionLevel = encodingCompressionLevel
	return c
}

func (c *Config) WithHttpClient(httpClient *http.Client) *Config {
	c.HttpClient = httpClient
	return c
}

func (c *Config) WithTlsConfig(tlsConfig *tls.Config) *Config {
	c.TlsConfig = tlsConfig
	return c
}

func (c *Config) WithRetryPolicy(maxRetryNumber int, retryStatus ...int) *Config {
	c.RetryStatus = retryStatus
	c.MaxRetries = maxRetryNumber
	return c
}

type Client struct {
	manager meilisearch.ServiceManager
	config  *Config
}

func New(ctx context.Context, config *Config) (*Client, error) {
	opt := make([]meilisearch.Option, 0)

	if config.ApiKey != "" {
		opt = append(opt, meilisearch.WithAPIKey(config.ApiKey))
	}

	if config.ContentEncoding != "" {
		opt = append(opt, meilisearch.WithContentEncoding(config.ContentEncoding, config.EncodingCompressionLevel))
	}

	if config.HttpClient != nil {
		opt = append(opt, meilisearch.WithCustomClient(config.HttpClient))
	}

	if config.TlsConfig != nil {
		opt = append(opt, meilisearch.WithCustomClientWithTLS(config.TlsConfig))
	}

	if config.MaxRetries > 0 {
		opt = append(opt, meilisearch.WithCustomRetries(config.RetryStatus, uint8(config.MaxRetries)))
	}

	sm := meilisearch.New(config.Host, opt...)

	context.AfterFunc(ctx, func() {
		sm.Close()
	})

	return &Client{
		manager: sm,
	}, nil
}

func (c *Client) Create(ctx context.Context) error {
	_, err := c.manager.CreateIndex(&meilisearch.IndexConfig{
		Uid:        c.config.CollectionName,
		PrimaryKey: "id",
	})

	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}

	return nil
}

type Data struct {
	Id      int          `json:"id"`
	Payload text.Payload `json:"payload"`
	Score   float64      `json:"_rankingScore,omitempty"`
}

func (c *Client) Store(ctx context.Context, id int, payload text.Payload) error {
	d := Data{
		Id:      id,
		Payload: payload,
	}

	idx := c.manager.Index(c.config.CollectionName)
	_, err := idx.UpdateDocuments([]Data{d}, "id")
	if err != nil {
		return fmt.Errorf("update documents: %w", err)
	}

	return nil
}

type SearchResult[T any] struct {
	Hits               []*T   `json:"hits"`
	Offset             int64  `json:"offset"`
	Limit              int64  `json:"limit"`
	EstimatedTotalHits int64  `json:"estimatedTotalHits"`
	ProcessingTimeMs   int64  `json:"processingTimeMs"`
	Query              string `json:"query"`
}

func (c *Client) Query(ctx context.Context, query string, option text.SearchOption) ([]text.Result, error) {
	idx := c.manager.Index(c.config.CollectionName)
	res, err := idx.SearchRaw(query, &meilisearch.SearchRequest{
		Limit:                 int64(option.Limit),
		Offset:                int64(option.Offset),
		Page:                  int64(option.Page),
		Sort:                  option.Sort,
		RankingScoreThreshold: option.ScoreThreshold,
	})
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	sr := &SearchResult[Data]{}
	if err := json.Unmarshal(*res, sr); err != nil {
		return nil, fmt.Errorf("unmarshal search result: %w", err)
	}

	results := make([]text.Result, 0, len(sr.Hits))
	for _, hit := range sr.Hits {
		results = append(results, text.Result{Id: hit.Id, Payload: hit.Payload, Score: hit.Score})
	}

	return results, nil
}

func (c *Client) Delete(ctx context.Context, id int) error {
	idx := c.manager.Index(c.config.CollectionName)

	_, err := idx.DeleteDocument(strconv.FormatInt(int64(id), 10))
	if err != nil {
		return fmt.Errorf("delete document: %w", err)
	}

	return nil
}

func (c *Client) Drop(ctx context.Context) error {
	_, err := c.manager.DeleteIndex(c.config.CollectionName)
	if err != nil {
		return fmt.Errorf("delete index: %w", err)
	}

	return nil
}

func (c *Client) UpdateSynonyms(ctx context.Context, synonyms map[string][]string) error {
	idx := c.manager.Index(c.config.CollectionName)

	_, err := idx.UpdateSynonyms(&synonyms)
	if err != nil {
		return fmt.Errorf("update synonyms: %w", err)
	}

	return nil
}
