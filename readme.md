# kafka practice in Go
概念手册请见 `./名词表.txt`

演示方法： `go test -v . -run 测试名`

比如 `go test -v . -run TestHowUseAsyncProducer0`
## TestHowUseAsyncProducer0
本例展示最简单的 异步生产者 的使用
## TestHowUseAsyncProducer1
本例展示 如何取得 异步生产者 的发送成功的返回
## TestHowUseAsyncProducer2
本例展示

0. 如何取得 异步生产者 的发送错误的返回
1. 限制发送的消息体的大小的上限

## TestHowUseBroker0
本例展示

0. 指定 topic 和 partition， 获得对应的 leader broker
1. 新建 topic

参考文章: [Creating Kafka topic in sarama - stackoverflow](https://stackoverflow.com/questions/44094926/creating-kafka-topic-in-sarama)

## TestHowUseConsumer0
本例展示最简单的 消费者 的使用
## TestHowUseConsumer1
本例展示最简单的 消费遗漏 的情景实现
## TestHowUseConsumer2
本例展示最简单的 重复消费 的情景实现
## TestHowUseConsumer3
本例展示最简单的 偏移量跟踪 的实现

参考文章：[Kafka入门（2）：消费与位移](https://zhuanlan.zhihu.com/p/163840121)
## TestHowUseConsumer4
本例展示 消费者组 的使用

参考文章：[sarama 源码学习 01：ConsumerGroup](https://zhuanlan.zhihu.com/p/109574627)
## TestHowUseSyncProducer
本例展示最简单的 同步生产者 的使用