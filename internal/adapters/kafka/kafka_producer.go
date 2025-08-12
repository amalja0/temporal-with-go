package kafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"sales-record-orchestration/internal/domain"
	"sales-record-orchestration/internal/ports"

	"github.com/IBM/sarama"
)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

func InitKafkaProducer(producer sarama.SyncProducer) ports.KafkaProducer {
	return &kafkaProducer{
		producer: producer,
	}
}

func (k *kafkaProducer) PublishSale(ctx context.Context, topic string, sale domain.SaleWithDetail) error {
	bytes, err := json.Marshal(sale)
	if err != nil {
		fmt.Println(err)
	}
	producer := k.producer
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(bytes),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

func (k *kafkaProducer) PublishSales(ctx context.Context, topic string, sales []domain.SaleWithDetail) error {
	for _, sale := range sales {
		k.PublishSale(ctx, topic, sale)
	}
	return nil
}
