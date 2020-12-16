package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

/*
	本例展示最简单的 消费者 的使用
*/

func HowUseConsumer0(stopCh chan bool, finishCh chan bool) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()
	partitionList, err := consumer.Partitions("test")
	if err != nil {
		panic(err)
	}
	for partition := range partitionList {
		pc, _ := consumer.ConsumePartition("test", int32(partition), sarama.OffsetNewest)
		defer pc.AsyncClose()
	L:
		for {
			select {
			case msg := <-pc.Messages():
				fmt.Printf("consumer<< Partition:%d, Offset:%d, Value: \"%s\"\n", msg.Partition, msg.Offset, string(msg.Value))
			case <-stopCh:
				break L
			}
		}
	}
	finishCh <- true
}

func TestHowUseConsumer0(t *testing.T) {
	stopCh, finishCh := make(chan bool), make(chan bool)
	go HowUseConsumer0(stopCh, finishCh)
	HowUseSyncProducer()
	stopCh <- true
	<- finishCh
}