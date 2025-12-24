package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/khusa-mahal/backend/internal/config"
	"github.com/khusa-mahal/backend/internal/models"
)

type SearchService struct {
	client *elasticsearch.Client
	index  string
}

func NewSearchService(cfg *config.Config) (*SearchService, error) {
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{cfg.Elasticsearch.URL},
	})
	if err != nil {
		return nil, err
	}

	return &SearchService{
		client: client,
		index:  cfg.Elasticsearch.Index,
	}, nil
}

// IndexProduct indexes a product for search
func (s *SearchService) IndexProduct(ctx context.Context, product *models.Product) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      s.index,
		DocumentID: product.ID.Hex(),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing product: %s", res.String())
	}

	return nil
}

// SearchProducts performs a full-text search
func (s *SearchService) SearchProducts(ctx context.Context, query string, from, size int) ([]models.Product, error) {
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"name^3", "description", "category^2"},
				"fuzziness": "AUTO",
			},
		},
		"from": from,
		"size": size,
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Hits []struct {
				Source models.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	products := make([]models.Product, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		products[i] = hit.Source
	}

	return products, nil
}

// DeleteProduct removes a product from the search index
func (s *SearchService) DeleteProduct(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      s.index,
		DocumentID: id,
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting product: %s", res.String())
	}

	return nil
}

// CreateIndex creates the products index with mappings
func (s *SearchService) CreateIndex(ctx context.Context) error {
	mapping := `{
		"mappings": {
			"properties": {
				"name": {
					"type": "text",
					"analyzer": "standard"
				},
				"description": {
					"type": "text"
				},
				"category": {
					"type": "keyword"
				},
				"price": {
					"type": "float"
				}
			}
		}
	}`

	res, err := s.client.Indices.Create(
		s.index,
		s.client.Indices.Create.WithContext(ctx),
		s.client.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		// Index might already exist, which is okay
		return nil
	}

	return nil
}
