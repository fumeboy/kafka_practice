package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

/*
	本例展示最简单的 同步生产者 的使用
*/

func HowUseSyncProducer() int64 {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	producer, _ := sarama.NewSyncProducerFromClient(client)
	defer producer.Close()

	topic := "test"
	var partition int32
	var offset int64
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("message %08d send", i)
		partition, offset, _ = producer.SendMessage(
			&sarama.ProducerMessage{
				Topic: topic,
				Key:   nil,
				Value: sarama.StringEncoder(text)})
		fmt.Println("Producer>> partition = ", partition, " offset = ", offset)
	}
	return offset
}

func TestHowUseSyncProducer(t *testing.T) {
	HowUseSyncProducer()
	ConsumeAll()
}