package main

import (
	"ProductRead/cache"
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

	redisClient, err := cache.NewRedisClientFromEnv()
	if err != nil {
		println("Warning: could not connect to Redis:", err)
		// Maybe i can add a retry mechanism here with a circuit breaker
	}
	if redisClient != nil {
		defer redisClient.Close()
	}

	repo := repositories.NewProductReadRepository(es, redisClient)

	r := gin.Default()

	r.GET("/products", handlers.GetAllProducts(repo))
	r.GET("/products/:id", handlers.GetProductById(repo))

	r.Run() // listens on :8080 by default
}
