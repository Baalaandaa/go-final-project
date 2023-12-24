package kafka

import (
	"context"
	"encoding/json"
	"final-project/internal/driver/model"
	"final-project/internal/driver/service"
	"github.com/juju/zaputil/zapctx"
	"github.com/segmentio/kafka-go"
	"strings"
	"time"
)

type adapter struct {
	config        *KafkaConfig
	driverService service.Driver
	consumer      *kafka.Reader
}

func (a adapter) Consume(ctx context.Context) error {
	a.consumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(a.config.ConsumeBroker, ","),
		Topic:          a.config.ConsumeTopic,
		SessionTimeout: 6 * time.Second,
	})
	defer a.consumer.Close()

	logger := zapctx.Logger(ctx).Sugar()

	for {
		msg, err := a.consumer.ReadMessage(ctx)
		if err != nil {
			logger.Errorf("Error consuming message: %+v", err)
			continue
		}
		var createdTripEvent model.CreatedTripEvent
		err = json.Unmarshal(msg.Value, &createdTripEvent)
		if err != nil {
			logger.Errorf("Can't unmarshal consumed message: %+v", err)
			continue
		}
		err = a.driverService.CreateTrip(ctx, createdTripEvent.Data)
		if err != nil {
			logger.Errorf("Can't process consumed message: %+v", err)
			continue
		}
		err = a.consumer.CommitMessages(ctx, msg)
		if err != nil {
			logger.Errorf("Error commiting message: %+v", err)
			continue
		}
	}
}

func (a adapter) Shutdown(ctx context.Context) {
	_ = a.consumer.Close()
}

func NewAdapter(config *KafkaConfig, driverService service.Driver) Adapter {
	return &adapter{
		driverService: driverService,
		config:        config,
	}
}
