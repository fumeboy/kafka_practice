package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
)

/*
	本例展示最简单的 重复消费 的情景实现
*/

func HowUseConsumer2(offset int64, length int) {
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
	i := 0
	for partition := range partitionList {
		pc, _ := consumer.ConsumePartition("test", int32(partition), offset-int64(length)+1)
		defer pc.AsyncClose()
		for ;i<length;i++{
			select {
			case msg := <-pc.Messages():
				fmt.Printf("consumer<< Partition:%d, Offset:%d, Value: \"%s\"\n", msg.Partition, msg.Offset, string(msg.Value))
			}
		}
	}
}


func TestHowUseConsumer2(t *testing.T) {
	offset := HowUseSyncProducer()
	HowUseConsumer2(offset, 5)
	HowUseConsumer2(offset, 5)
}

/*
	解释：
	该例运行后
		0 生产者，先生产 5 条消息, 返回最新的消息的偏移量，命名为 offset
		1 消费者，标记 offset-5+1 为起始偏移值， 成功消费了偏移量是 offset-5+1 到 offset 之间的这些消息
		2 再启动一个消费者，同上一步骤，因为起始偏移量没有更新，所以重复消费了这些消息
*/