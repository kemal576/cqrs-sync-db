package consumers

import (
	"encoding/json"
	"log"
	"strings"

	"ProductRead/models"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
)

type ProductConsumer struct {
	KafkaBroker string
	Topic       string
	EsClient    *elasticsearch.Client
}

func NewProductConsumer(kafkaBroker, topic, esAddr string) (*ProductConsumer, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{esAddr},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ProductConsumer{
		KafkaBroker: kafkaBroker,
		Topic:       topic,
		EsClient:    es,
	}, nil
}

const (
	OpCreate = "c"
	OpRead   = "r"
	OpUpdate = "u"
	OpDelete = "d"
)

type DebeziumMessage struct {
	Op     string          `json:"op"`
	After  json.RawMessage `json:"after"`
	Before json.RawMessage `json:"before"`
}

func (pc *ProductConsumer) Consume() {
	consumer, err := sarama.NewConsumer([]string{pc.KafkaBroker}, nil)
	if err != nil {
		log.Fatalf("Kafka error: %v", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(pc.Topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Partition error: %v", err)
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		var dbzMsg DebeziumMessage
		if err := json.Unmarshal(msg.Value, &dbzMsg); err != nil {
			log.Printf("Invalid CDC message: %v", err)
			continue
		}

		switch dbzMsg.Op {
		case OpCreate, OpRead, OpUpdate:
			pc.handleUpsert(dbzMsg.After)
		case OpDelete:
			pc.handleDelete(dbzMsg.Before)
		default:
			log.Printf("Unknown op type: %s", dbzMsg.Op)
		}
	}
}

func (pc *ProductConsumer) handleUpsert(raw json.RawMessage) {
	var dataStr string
	if err := json.Unmarshal(raw, &dataStr); err != nil {
		log.Printf("after field is not a string: %v", err)
		return
	}
	var product models.Product
	if err := json.Unmarshal([]byte(dataStr), &product); err != nil {
		log.Printf("Invalid product in after: %v", err)
		return
	}
	_, err := pc.EsClient.Index(
		"products",
		strings.NewReader(dataStr),
		pc.EsClient.Index.WithDocumentID(product.Id.String()),
		pc.EsClient.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Printf("Elasticsearch write error: %v", err)
	} else {
		log.Printf("Product %s indexed", product.Id)
	}
}

func (pc *ProductConsumer) handleDelete(raw json.RawMessage) {
	var dataStr string
	if err := json.Unmarshal(raw, &dataStr); err != nil {
		log.Printf("before field is not a string: %v", err)
		return
	}
	var product models.Product
	if err := json.Unmarshal([]byte(dataStr), &product); err != nil {
		log.Printf("Invalid product in before: %v", err)
		return
	}
	_, err := pc.EsClient.Delete("products", product.Id.String())
	if err != nil {
		log.Printf("Elasticsearch delete error: %v", err)
	} else {
		log.Printf("Product %s deleted", product.Id)
	}
}
