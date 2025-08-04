package main

import (
	"ProductConsumer/consumers"
	"log"
	"os"
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
	c.Consume()
}
