package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
	"time"
)

/*
	本例展示最简单的 偏移量跟踪 的实现
*/

func HowUseConsumer3(stopCh chan bool, finishCh chan bool) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 自动 commit offset 时间间隔
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	consumer, _ := sarama.NewConsumerFromClient(client)
	offsetManager,_:=sarama.NewOffsetManagerFromClient("c1",client) // 偏移量管理器
	defer offsetManager.Close()
	partitionList, _ := consumer.Partitions("test")
	for _, partition := range partitionList {
		partitionOffsetManager,_:=offsetManager.ManagePartition("test",partition) // 对应分区的偏移量管理器
		nextOffset,_:=partitionOffsetManager.NextOffset() // 取得下一消息的偏移量
		pc, _ := consumer.ConsumePartition("test", partition, nextOffset)
		defer pc.AsyncClose()
		defer partitionOffsetManager.Close()

		var i int64 = 0
	L:
		for {
			select {
			case msg := <-pc.Messages():
				fmt.Printf("consumer<< Partition:%d, Offset:%d, Value: \"%s\"\n", msg.Partition, msg.Offset, string(msg.Value))
				i++
			case <- stopCh:
				partitionOffsetManager.MarkOffset(nextOffset+i,"modified metadata")
				// MarkOffset 更新最后消费的 offset
				offsetManager.Commit() // 手动 commit 上句代码执行后的 Marked Offset
				break L
			}
		}
	}
	finishCh <- true
}

func TestHowUseConsumer3(t *testing.T) {
	HowUseSyncProducer()
	stopCh, finishCh := make(chan bool), make(chan bool)
	go HowUseConsumer3(stopCh, finishCh)
	HowUseSyncProducer()
	time.Sleep(10*time.Millisecond)
	stopCh <- true
	<- finishCh
}

/*
	和 TestHowUseConsumer1 对比
	会发现这时候没有了之前 消费遗漏 的问题，这是因为 kafka 可以保存每个分区的消费的 offset 值
	这样我们只需要读这个偏移值，从这个值继续就可以了，这就是 偏移量跟踪 offset tracking
	结束前，需要 commit 新的偏移量
	sarama 有 auto commit 机制，但结束前还是应该 手动 commit 一次
*/