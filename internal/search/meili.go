package search

import (
	"context"

	"github.com/meilisearch/meilisearch-go"
)

type Client struct {
	cli *meilisearch.Client
}

func New(url, key string) *Client {
	return &Client{
		cli: meilisearch.NewClient(meilisearch.ClientConfig{
			Host:   url,
			APIKey: key,
		}),
	}
}

// IndexPoints indexes/updates map points in batch.
func (c *Client) IndexPoints(ctx context.Context, docs any) error {
	_, err := c.cli.Index("map_points").AddDocuments(docs, "id")
	return err
}

// QueryPoints supports autocomplete & filters.
func (c *Client) QueryPoints(ctx context.Context, q string, limit int) (any, error) {
	return c.cli.Index("map_points").Search(q, &meilisearch.SearchRequest{
		Limit: int64(limit),
	})
}
