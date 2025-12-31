package consumers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"ProductConsumer/models"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
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

type DebeziumMessageValue struct {
	Payload struct {
		Op     string          `json:"op"`
		After  json.RawMessage `json:"after"`
		Before json.RawMessage `json:"before"`
	} `json:"payload"`
}

type DebeziumMessageKey struct {
	Id uuid.UUID `json:"id"`
}

func (pc *ProductConsumer) Consume() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{pc.KafkaBroker}, config)
	if err != nil {
		log.Printf("Kafka connection error: %v. Retrying...", err)
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(pc.Topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Printf("Partition error: %v. Retrying...", err)
		return
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		if len(msg.Value) == 0 {
			log.Printf("Skipping tombstone message for key: %s", string(msg.Key))
			continue
		}

		log.Printf("Raw Kafka message value: %s", string(msg.Value))

		var dbzMsgValue DebeziumMessageValue
		if err := json.Unmarshal(msg.Value, &dbzMsgValue); err != nil {
			log.Printf("Invalid CDC message value: %v", err)
			continue
		}

		switch dbzMsgValue.Payload.Op {
		case OpCreate, OpRead, OpUpdate:
			pc.handleUpsert(dbzMsgValue.Payload.After)

		case OpDelete:
			pc.handleDelete(msg.Key)

		default:
			log.Printf("Unknown operation type: %s", dbzMsgValue.Payload.Op)
			log.Printf("Kafka Message: Key: %s, Value: %s", string(msg.Key), string(msg.Value))
		}
	}
}

func (pc *ProductConsumer) handleUpsert(raw json.RawMessage) {
	var dataStr string
	if err := json.Unmarshal(raw, &dataStr); err != nil {
		log.Printf("after field is not a string: %v", err)
		return
	}

	log.Printf("Received upsert data: %s", dataStr)

	var product models.Product
	if err := product.UnmarshalJSON([]byte(dataStr)); err != nil {
		log.Printf("Invalid product in after: %v", err)
		return
	}

	productJson, err := json.Marshal(product)
	if err != nil {
		log.Printf("Error marshaling product: %v", err)
		return
	}

	res, err := pc.EsClient.Index(
		"products",
		bytes.NewReader(productJson),
		pc.EsClient.Index.WithDocumentID(product.Id.String()),
		pc.EsClient.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Printf("Elasticsearch write error: %v", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("Elasticsearch index error: %s", string(bodyBytes))
		return
	}

	log.Printf("Product %s indexed", product.Id)
}

func (pc *ProductConsumer) handleDelete(key []byte) {
	log.Printf("Deleting product with ID: %s", string(key))

	var keyModel DebeziumMessageKey
	if err := json.Unmarshal(key, &keyModel); err != nil {
		log.Printf("Failed to parse key payload: %v", err)
		return
	}

	res, err := pc.EsClient.Delete("products", keyModel.Id.String())
	if err != nil {
		log.Printf("Elasticsearch delete error: %v", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		log.Printf("Elasticsearch delete error: %s", string(bodyBytes))
		return
	}

	log.Printf("Product %s deleted", keyModel.Id)
}
