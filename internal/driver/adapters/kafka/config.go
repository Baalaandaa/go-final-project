package kafka

type KafkaConfig struct {
	ConsumeBroker string `env:"KAFKA_CONSUME_BROKER"`
	ProduceBroker string `env:"KAFKA_PRODUCE_BROKER"` // divided by ','
	ConsumeTopic  string `env:"KAFKA_CONSUME_TOPIC"`
}
