package main

import (
	"ProductRead/consumers"
	"ProductRead/handlers"
	"ProductRead/repositories"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

func main() {
	esAddr := os.Getenv("ELASTICSEARCH_ADDR")
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

	consumer, err := consumers.NewProductConsumer(os.Getenv("KAFKA_BROKER"), os.Getenv("PRODUCT_CDC_TOPIC"), esAddr)
	if err != nil {
		panic(err)
	}
	go consumer.Consume()
}
