package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/guobinqiu/kafka-demo/test-exactly-once/internal/message"
	_ "github.com/mattn/go-sqlite3"
)

// Kafka 配置
var config = &kafka.ConfigMap{
	"bootstrap.servers": "127.0.0.1:29092,127.0.0.1:39092,127.0.0.1:49092", // Kafka 集群地址（多个 broker 提高容错能力）
	"group.id":          "go-consumer-group",                               // 消费者所属的消费组 ID，Kafka 用这个来做分区协调、偏移量管理等。
	"auto.offset.reset": "latest",                                          //"latest" 表示从最新的消息开始消费；"earliest" 表示从最旧的消息开始消费。
	"security.protocol": "PLAINTEXT",                                       // 安全协议类型。 "PLAINTEXT" 表示不加密、无认证的明文通信；还可以设置为 "SSL"、"SASL_PLAINTEXT"、"SASL_SSL" 等
	// "isolation.level":    "read_committed",                                  // 如果使用了事务生产者，并希望只读取已提交的消息
	"enable.auto.commit": false, // 禁用自动提交
}

func main() {
	// 初始化 SQLite 数据库
	db, err := sql.Open("sqlite3", "./consumer_state.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 创建幂等性校验表（如果不存在）
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS consumed_messages (
		id TEXT PRIMARY KEY
	);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		log.Fatalf("Error creating consumer: %v", err)
	}
	defer consumer.Close()

	topicName := "test"
	err = consumer.Subscribe(topicName, nil)
	if err != nil {
		log.Fatalf("Error Subscribe failed: %v", err)
	}

	log.Println("Consumer started...")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Consumer exiting...")
				done <- struct{}{}
				return
				// default:
				// msg, err := consumer.ReadMessage(-1) // 阻塞直到消息到达 这样永远没有机会执行到 case <-ctx.Done()
				// if err != nil {
				// 	log.Printf("ReadMessage error: %v", err)
				// 	continue
				// }

				//event := consumer.Poll(100) // 非阻塞轮询消息 还是没有走到 case <-ctx.Done() 还是阻塞的 阻塞100ms
			case <-time.After(100 * time.Millisecond): // 为防止CPU空转 自定义轮询间隔
				event := consumer.Poll(0) // 完全不阻塞
				if event == nil {
					continue
				}

				switch msg := event.(type) {
				case *kafka.Message:
					var m message.Message
					if err := json.Unmarshal(msg.Value, &m); err != nil {
						log.Printf("JSON error: %v", err)
						continue
					}

					// ID 作为唯一标识
					if m.ID == "" {
						log.Println("Empty message ID, skipped")
						continue
					}

					// 幂等校验
					var existingID string
					err := db.QueryRow(`SELECT id FROM consumed_messages WHERE id = ?`, m.ID).Scan(&existingID)
					if err != nil && err != sql.ErrNoRows {
						log.Printf("DB query error: %v", err)
						continue
					}
					if existingID != "" {
						log.Printf("Duplicate message skipped (ID=%s, content=%s)", m.ID, m.Content)
						continue
					}

					// 开启事务
					tx, _ := db.Begin()

					// 处理消息
					log.Printf("Consumed message: %s", m.Content)
					time.Sleep(time.Second)

					// 处理成功后再写入幂等表
					_, err = tx.Exec(`INSERT INTO consumed_messages (id) VALUES (?)`, m.ID)
					if err != nil {
						log.Printf("DB insert error: %v", err)
						_ = tx.Rollback()
						continue
					}

					// 成功处理消息后提交偏移量
					if _, err := consumer.CommitMessage(msg); err != nil {
						log.Printf("Commit error: %v", err)
						_ = tx.Rollback()
						continue
					}

					// 提交事务
					_ = tx.Commit()
				default:
					// 忽略其他事件
				}
			}
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // 阻塞直到收到退出信号
	log.Println("Shutting down...")
	cancel() // 触发ctx.Done()

	<-done
	log.Println("Consumer exited.")
}
