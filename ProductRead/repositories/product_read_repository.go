package repositories

import (
	"ProductRead/cache"
	"ProductRead/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

type ProductReadRepository struct {
	es    *elasticsearch.Client
	cache *cache.RedisClient
}

func NewProductReadRepository(es *elasticsearch.Client, c *cache.RedisClient) *ProductReadRepository {
	return &ProductReadRepository{es: es, cache: c}
}

func (r *ProductReadRepository) GetById(ctx context.Context, id string) (*models.Product, error) {
	if r.cache != nil {
		if s, err := r.cache.Get(ctx, "product:"+id); err == nil && s != "" {
			var p models.Product
			if err := json.Unmarshal([]byte(s), &p); err == nil {
				return &p, nil
			}
		}
	}
	res, err := r.es.Get("products", id)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("error getting product: %s", res.String())
	}

	var doc struct {
		Source models.Product `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&doc); err != nil {
		return nil, err
	}

	if r.cache != nil {
		if b, err := json.Marshal(doc.Source); err == nil {
			// set with 60s TTL to test easily, can be read from environment variables one day
			_ = r.cache.Set(ctx, "product:"+id, string(b), 60*time.Second)
		}
	}

	return &doc.Source, nil
}

func (r *ProductReadRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	var products []models.Product

	query := `{"query": {"match_all": {}}}`

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("products"),
		r.es.Search.WithBody(strings.NewReader(query)),
		r.es.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.IsError() {
		return nil, fmt.Errorf("error searching products: %s", res.String())
	}

	var resp struct {
		Hits struct {
			Hits []struct {
				Source models.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	for _, hit := range resp.Hits.Hits {
		products = append(products, hit.Source)
	}

	return products, nil
}
