package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"time"
)

func ConsumeAll() {
	config := sarama.NewConfig()
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	consumer, _ := sarama.NewConsumerFromClient(client)
	offsetManager,_:=sarama.NewOffsetManagerFromClient("g1",client)
	defer offsetManager.Close()
	partitionList, _ := consumer.Partitions("test")
	for _, partition := range partitionList {
		partitionOffsetManager,_:=offsetManager.ManagePartition("test",partition)
		nextOffset,_:=partitionOffsetManager.NextOffset()
		pc, _ := consumer.ConsumePartition("test", partition, nextOffset)
		defer pc.AsyncClose()
		defer partitionOffsetManager.Close()

		timer:=time.NewTimer(time.Millisecond*10)
		var i int64 = 0
	L:
		for {
			select {
			case msg := <-pc.Messages():
				timer.Reset(time.Millisecond*1)
				fmt.Printf("consumer<< Partition:%d, Offset:%d, Value: \"%s\"\n", msg.Partition, msg.Offset, string(msg.Value))
				i++
			case <- timer.C:
				partitionOffsetManager.MarkOffset(nextOffset+i,"modified metadata")
				break L
			}
		}
	}
}
