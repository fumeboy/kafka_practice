package kafka_practice

import (
	"fmt"
	"github.com/Shopify/sarama"
	"testing"
	"time"
)

/*
	本例展示
		0 指定 topic 和 partition， 获得对应的 leader broker
		1 新建 topic
*/

func HowUseBroker0() {
	config := sarama.NewConfig()
	client, err := sarama.NewClient([]string{"localhost:9092"}, config)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	leaderBroker,_ := client.Leader("test", 0)

	topicDetail := &sarama.TopicDetail{}
	topicDetail.NumPartitions = int32(1) // 分区数
	topicDetail.ReplicationFactor = int16(1) // 备份数
	topicDetail.ConfigEntries = make(map[string]*string) // 不知道
	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails["new_topic_from_client_test"] = topicDetail
	resp, err := leaderBroker.CreateTopics(&sarama.CreateTopicsRequest{
		TopicDetails: topicDetails,
		Timeout:      time.Second * 15,
	})
	if err != nil {
		panic(err)
	}
	t := resp.TopicErrors
	for key, val := range t {
		fmt.Printf("Key: '%s', Err: %#v, ErrMsg: %#v\n", key, val.Err.Error(), val.ErrMsg)
	}
	client.RefreshMetadata() // 重新获取元数据
	fmt.Println(client.Topics())
}

func TestHowUseBroker0(t *testing.T) {
	HowUseBroker0()
}