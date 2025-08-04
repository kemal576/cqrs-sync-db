package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
}

// Yardımcı struct
type rawProduct struct {
	Id        uuid.UUID `json:"_id"`
	CreatedAt struct {
		Date int64 `json:"$date"`
	} `json:"createdAt"`
	UpdatedAt struct {
		Date int64 `json:"$date"`
	} `json:"updatedAt"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       struct {
		NumberDecimal string `json:"$numberDecimal"`
	} `json:"price"`
}

func (p *Product) UnmarshalJSON(data []byte) error {
	var raw rawProduct
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	price, err := strconv.ParseFloat(raw.Price.NumberDecimal, 64)
	if err != nil {
		return fmt.Errorf("price parse error: %w", err)
	}

	p.Id = raw.Id
	p.Name = raw.Name
	p.Description = raw.Description
	p.Price = price
	p.CreatedAt = time.UnixMilli(raw.CreatedAt.Date)
	p.UpdatedAt = time.UnixMilli(raw.UpdatedAt.Date)

	return nil
}
