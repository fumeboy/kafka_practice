package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

/*
	本例展示
		如何取得 异步生产者 的发送错误的返回
		限制发送的消息体的大小的上限
*/

func HowUseAsyncProducer2() {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = true // 返回错误信息
	config.Producer.MaxMessageBytes = 1 // 令消息体最大大小为 1 byte，用以发生错误
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	producer, _ := sarama.NewAsyncProducerFromClient(client)
	defer producer.Close()

	fmt.Println("Producer Input() cap ==", cap(producer.Input()))
	topic := "test"
	for i := 0; i < 5; i++ {
		producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Key:   nil,
			Value: sarama.StringEncoder("01"),
		}
		r := <- producer.Errors()
		// 取出错误信息, [!important] 如果不取出，将阻塞 producer.Input()
		fmt.Println("Producer>> Err = ", r.Err)
	}
}

func TestHowUseAsyncProducer2(t *testing.T) {
	HowUseAsyncProducer2()
}

