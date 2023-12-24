package kafka_producer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type producer struct {
	brokers []string

	producers map[string]*kafka.Writer
}

func (p producer) getProducer(topicName string) *kafka.Writer {
	if writer, ok := p.producers[topicName]; ok {
		return writer
	}
	p.producers[topicName] = &kafka.Writer{
		Addr:  kafka.TCP(p.brokers...),
		Topic: topicName,
	}
	return p.producers[topicName]
}

func (p producer) Produce(ctx context.Context, message *KafkaMessage, topicName string) error {
	writer := p.getProducer(topicName)
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return writer.WriteMessages(ctx, kafka.Message{Key: []byte(message.ID), Value: messageBytes})
}

func NewProducer(brokers []string) Producer {
	return &producer{
		brokers:   brokers,
		producers: make(map[string]*kafka.Writer),
	}
}
