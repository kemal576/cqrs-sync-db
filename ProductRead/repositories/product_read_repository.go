package repositories

import (
	"ProductRead/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type ProductReadRepository struct {
	es *elasticsearch.Client
}

func NewProductReadRepository(es *elasticsearch.Client) *ProductReadRepository {
	return &ProductReadRepository{es: es}
}

func (r *ProductReadRepository) GetById(ctx context.Context, id string) (*models.Product, error) {
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
