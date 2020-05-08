package eventqueue

import (
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

// KafkaConfig the config to kafka
type KafkaConfig struct {
	Broker string
}

type kafkaProducer struct {
	producer     *kafka.Producer
	deliveryChan chan kafka.Event
}

func (k *kafkaProducer) Produce(topic string, message string) error {
	k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, k.deliveryChan)

	e := <-k.deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}
	return nil
}

func (k *kafkaProducer) Close() {
	close(k.deliveryChan)
}

// NewProducer constructor function to kafka producer
func NewProducer(config *KafkaConfig) (Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.Broker})
	if err != nil {
		return nil, err
	}

	k := kafkaProducer{p, make(chan kafka.Event)}
	return &k, nil
}
