package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

/*
	本例展示最简单的 异步生产者 的使用
*/

func HowUseAsyncProducer0() {
	config := sarama.NewConfig()
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	topic := "test"
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("message %08d send", i)
		producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Key:   nil,
			Value: sarama.StringEncoder(text),
		}
	}
}

func TestHowUseAsyncProducer0(t *testing.T) {
	HowUseAsyncProducer0()
	ConsumeAll()
}
