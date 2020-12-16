package kafka_practice

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"sync"
	"testing"
	"time"
)

/*
	本例展示 消费者组 的使用
*/

func HowUseConsumer4(stopCh chan bool, group string, name string) {
	config := sarama.NewConfig()
	client, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, group, config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	handler := consumerGroupHandler{
		name: name,
		ready: make(chan bool),
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			/*
				需要注意每当遇到 Rebalance
				需要再次执行 Consume() 来恢复连接
				应当在一个无限循环中不停地调用 Consume()
			*/
			if err := client.Consume(ctx, []string{"test"}, &handler); err != nil {
				panic(err)
			}
			// 如果 context 被 cancel 了，那么退出
			if ctx.Err() != nil {
				return
			}
			handler.name = name
			handler.ready = make(chan bool)
		}
	}()
	<-handler.ready // 等待 Consume() 执行完毕
	fmt.Println("consumer up and running!...")
	<- stopCh
	cancel()
	wg.Wait()
}


type consumerGroupHandler struct {
	name string
	ready chan bool
}

// Setup 执行在 获得新 session 后 的第一步, 在 ConsumeClaim() 之前
func (consumer *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready

	close(consumer.ready)
	return nil
}

// Cleanup 执行在 session 结束前, 当所有 ConsumeClaim goroutines 都退出时
func (consumer *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// ConsumeClaim() 必须启动一个消费循环
	// claim 和 session 的含义详见 `./名词表.txt`
	for message := range claim.Messages() {
		fmt.Printf("%s consumed: value = %s, offset = %v, topic = %s\n", consumer.name, string(message.Value), message.Offset, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}


func TestHowUseConsumer4(t *testing.T) {
	StopCh,StopCh1,StopCh2 := make(chan bool),make(chan bool),make(chan bool)

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		HowUseConsumer4(StopCh, "g1", "c1")
		wg.Done()
	}()
	go func() {
		time.Sleep(10*time.Millisecond)
		HowUseConsumer4(StopCh1, "g1", "c2") // 观察是否会重复消费 c1 已经消费过的消息
		wg.Done()
	}()
	go func() {
		time.Sleep(20*time.Millisecond)
		HowUseConsumer4(StopCh2, "g2", "c3") // 和上面两个消费者不在同一组
		wg.Done()
	}()

	HowUseSyncProducer()
	time.Sleep(1000*time.Millisecond)
	StopCh <- true
	StopCh1 <- true
	StopCh2 <- true

	wg.Wait()
}