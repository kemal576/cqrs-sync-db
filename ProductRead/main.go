package main

import (
	"ProductRead/handlers"
	"ProductRead/repositories"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

func main() {
	esAddr := os.Getenv("ELASTICSEARCH_ADDR")
	println("Elasticsearch address:", esAddr)
	if esAddr == "" {
		panic("ELASTICSEARCH_ADDR environment variable is not set")
	}

	cfg := elasticsearch.Config{
		Addresses: []string{esAddr},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	repo := repositories.NewProductReadRepository(es)

	r := gin.Default()

	r.GET("/products", handlers.GetAllProducts(repo))
	r.GET("/products/:id", handlers.GetProductById(repo))

	r.Run() // listens on :8080 by default
}
