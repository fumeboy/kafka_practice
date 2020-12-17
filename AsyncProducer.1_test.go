package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

/*
	本例展示 如何取得 异步生产者 的发送成功的返回
*/

func HowUseAsyncProducer1() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true // 返回成功信息
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	producer, _ := sarama.NewAsyncProducerFromClient(client)
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
	for i := 0; i < 5; i++ {
		r := <- producer.Successes()
		fmt.Println("Producer>> partition = ", r.Partition, " offset = ", r.Offset)
	}
}

func TestHowUseAsyncProducer1(t *testing.T) {
	HowUseAsyncProducer1()
	ConsumeAll()
}

