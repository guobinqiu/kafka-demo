package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/guobinqiu/kafka-demo/test-consumer-parallel/internal/message"
)

const (
	numPartitions     = 3 // 分区数
	replicationFactor = 1 // 副本数小于等于broker数量
	numMessages       = 9
)

var topicName = "test-consumer-parallel"

func main() {
	// Kafka 配置
	config := &kafka.ConfigMap{
		"bootstrap.servers":   "127.0.0.1:29092,127.0.0.1:39092,127.0.0.1:49092", // Kafka 集群地址（多个 broker 提高容错能力）
		"acks":                "all",                                             // all 表示所有 ISR（In-Sync Replicas）都确认后才认为发送成功，确保数据可靠性
		"enable.idempotence":  "true",                                            // 开启幂等性
		"retries":             5,                                                 // 重试次数（发送失败后重试），适用于短暂的错误（如网络闪断、leader 切换等瞬时问题）
		"retry.backoff.ms":    1000,                                              // 两次重试之间的时间间隔（毫秒）
		"delivery.timeout.ms": 30000,                                             // 发送一条消息的总时间限制 包含重试次数时间
		// "transactional.id":    "t1",                                              // 启用事务，id要唯一
	}

	// 创建 Kafka AdminClient 用于创建 topic
	adminClient, err := kafka.NewAdminClient(config)
	if err != nil {
		log.Fatalf("Error creating admin client: %v", err)
	}
	defer adminClient.Close()

	ctx := context.TODO()

	createTopic(ctx, adminClient, topicName, numPartitions, replicationFactor)

	producer, err := kafka.NewProducer(config)
	if err != nil {
		log.Fatalf("Error creating producer: %v", err)
	}
	defer producer.Close()

	// err = producer.InitTransactions(ctx)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize transactions: %v", err)
	// }

	// 使用 deliveryChan 接收发送消息的结果
	deliveryChan := make(chan kafka.Event, numMessages)

	var wg sync.WaitGroup

	for i := 0; i < numMessages; i++ {
		message := message.Message{
			ID: fmt.Sprintf("%d", i+1), //消息唯一ID, 给消费端做幂等用
			//OrderID: fmt.Sprintf("%d", i+1),          //订单ID, 按Key分区用
			Content: fmt.Sprintf("Message #%d", i+1), //消息内容
		}

		value, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshalling message: %v", err)
			continue
		}

		// _ = producer.BeginTransaction()

		// 发送消息 (异步)
		err = producer.Produce(&kafka.Message{
			// TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: kafka.PartitionAny},
			TopicPartition: kafka.TopicPartition{Topic: &topicName, Partition: int32(i % numPartitions)}, // 为了消费端负载均衡自己实现分区轮询
			Value:          value,
			//Key:            []byte(message.OrderID), // 这里指定 key, 用于保证同一个 OrderID 落入同一分区, 感觉Key分区不是很均衡, 算法为: hash(key) % numPartitions
		}, deliveryChan) // 没有调用 Flush(timeout) 就必须手动监控 DeliveryChan 来确保消息投递成功

		if err != nil {
			log.Printf("Error producing message: %v", err)
			// _ = producer.AbortTransaction(ctx)
			continue
		}
		wg.Add(1)
	}

	// 监听消息的发送结果
	go func() {
		for msg := range deliveryChan {
			m := msg.(*kafka.Message)
			if m.TopicPartition.Error != nil {
				// 发送失败，打印消息并回滚事务
				// 这里你可以发送到死信队列 或做业务层面重试逻辑 或其他兜底机制
				log.Printf("Delivery failed (topic: %s, partition: %d, offset: %d, error: %v)",
					*m.TopicPartition.Topic,
					m.TopicPartition.Partition,
					m.TopicPartition.Offset,
					m.TopicPartition.Error,
				)
				// _ = producer.AbortTransaction(ctx)
			} else {
				// 发送成功，打印消息并提交事务
				log.Printf("Message delivered (topic: %s, partition: %d, offset: %d, key: %s, value: %s)",
					*m.TopicPartition.Topic,
					m.TopicPartition.Partition,
					m.TopicPartition.Offset,
					string(m.Key),
					string(m.Value),
				)
				// _ = producer.CommitTransaction(ctx)
			}
			wg.Done()
		}
	}()

	wg.Wait()
	close(deliveryChan)
}

func createTopic(ctx context.Context, adminClient *kafka.AdminClient, topic string, numPartitions int, replicationFactor int) {
	results, err := adminClient.CreateTopics(
		ctx,
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     numPartitions,
			ReplicationFactor: replicationFactor,
			Config: map[string]string{
				"retention.ms":    "604800000",  // 设置保留时间为7天
				"retention.bytes": "1073741824", // 设置最大保留大小为1GB
			},
		}},
		kafka.SetAdminOperationTimeout(10*time.Second),
	)

	if err != nil {
		log.Fatal("Kafka CreateTopics request failed:", err)
	}

	for _, result := range results {
		if result.Error.Code() == kafka.ErrNoError {
			log.Printf("Successfully created topic %s.\n", result.Topic)
		} else if result.Error.Code() == kafka.ErrTopicAlreadyExists {
			log.Printf("Topic %s already exists.\n", result.Topic)
		} else {
			log.Fatalf("Failed to create topic %s: %v\n", result.Topic, result.Error)
		}
	}
}
