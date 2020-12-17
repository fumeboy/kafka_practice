package kafka_practice

import (
	"testing"
)

/*
	本例展示最简单的 消费遗漏 的情景实现
*/

func TestHowUseConsumer1(t *testing.T) {
	HowUseSyncProducer()
	stopCh, finishCh := make(chan bool), make(chan bool)
	go HowUseConsumer0(stopCh, finishCh)
	HowUseSyncProducer()
	stopCh <- true
	<- finishCh
}

/*
	解释：
	该例运行后
		0 生产者，先生产 5 条消息
		1 消费者，标记 sarama.OffsetNewest 为起始偏移值，读指针移动到队列末尾，因此直接跳过了步骤 0 生产的消息
		2 生产者，再生产 5 条消息
		3 消费者读到了起始偏移值之后的消息，成功消费这 5 条消息
*/