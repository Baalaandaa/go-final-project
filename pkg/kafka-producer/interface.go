package kafka_producer

import (
	"context"
)

type KafkaMessage struct {
	ID              string      `json:"id"`
	Source          string      `json:"source"`
	Type            string      `json:"type"`
	DataContentType string      `json:"datacontenttype"`
	Time            string      `json:"time"`
	Data            interface{} `json:"data"`
}

type Producer interface {
	Produce(ctx context.Context, message *KafkaMessage, topicName string) error
}
