package main

import (
	"ProductConsumer/consumers"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("PRODUCT_CDC_TOPIC")
	esAddr := os.Getenv("ELASTICSEARCH_ADDR")

	if kafkaBroker == "" || topic == "" || esAddr == "" {
		log.Fatal("KAFKA_BROKER, PRODUCT_CDC_TOPIC, and ELASTICSEARCH_ADDR env variables must be set")
	}

	c, err := consumers.NewProductConsumer(kafkaBroker, topic, esAddr)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	// Start Kafka consumer in a goroutine with retry loop
	go func() {
		for {
			c.Consume()
			log.Println("Kafka consumer disconnected, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}()

	// Start HTTP server for health checks
	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
