# Kafka Demo

## docker 安装

> https://github.com/guobinqiu/kafka-demo/tree/main/install/docker

运行

```
docker compose -f install/docker/xxx.yaml > up -d
```

### apache kafka多实例 基于kraft

```
services:
  kafka-0:
    container_name: kafka-0
    image: apache/kafka:3.9.0
    ports:
      - 29092:29092
    environment:
      KAFKA_NODE_ID: 0
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_CONTROLLER_QUORUM_VOTERS: 0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:29092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-0:9092,CONTROLLER://kafka-0:9093,EXTERNAL://localhost:29092
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT

  kafka-1:
    container_name: kafka-1
    image: apache/kafka:3.9.0
    ports:
      - 39092:39092
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_CONTROLLER_QUORUM_VOTERS: 0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:39092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-1:9092,CONTROLLER://kafka-1:9093,EXTERNAL://localhost:39092
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT

  kafka-2:
    container_name: kafka-2
    image: apache/kafka:3.9.0
    ports:
      - 49092:49092
    environment:
      KAFKA_NODE_ID: 2
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_CONTROLLER_QUORUM_VOTERS: 0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:49092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka-2:9092,CONTROLLER://kafka-2:9093,EXTERNAL://localhost:49092
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
```

### apache kafka多实例 基于kraft [官方配置]

更加高可用但也需要更多节点
 
```
services:
  controller-1:
    image: apache/kafka:latest
    container_name: controller-1
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_PROCESS_ROLES: controller
      KAFKA_LISTENERS: CONTROLLER://:9093
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@controller-1:9093,2@controller-2:9093,3@controller-3:9093
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0

  controller-2:
    image: apache/kafka:latest
    container_name: controller-2
    environment:
      KAFKA_NODE_ID: 2
      KAFKA_PROCESS_ROLES: controller
      KAFKA_LISTENERS: CONTROLLER://:9093
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@controller-1:9093,2@controller-2:9093,3@controller-3:9093
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0

  controller-3:
    image: apache/kafka:latest
    container_name: controller-3
    environment:
      KAFKA_NODE_ID: 3
      KAFKA_PROCESS_ROLES: controller
      KAFKA_LISTENERS: CONTROLLER://:9093
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@controller-1:9093,2@controller-2:9093,3@controller-3:9093
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0

  broker-1:
    image: apache/kafka:latest
    container_name: broker-1
    ports:
      - 29092:9092
    environment:
      KAFKA_NODE_ID: 4
      KAFKA_PROCESS_ROLES: broker
      KAFKA_LISTENERS: 'PLAINTEXT://:19092,PLAINTEXT_HOST://:9092'
      KAFKA_ADVERTISED_LISTENERS: 'PLAINTEXT://broker-1:19092,PLAINTEXT_HOST://localhost:29092'
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@controller-1:9093,2@controller-2:9093,3@controller-3:9093
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
    depends_on:
      - controller-1
      - controller-2
      - controller-3

  broker-2:
    image: apache/kafka:latest
    container_name: broker-2
    ports:
      - 39092:9092
    environment:
      KAFKA_NODE_ID: 5
      KAFKA_PROCESS_ROLES: broker
      KAFKA_LISTENERS: 'PLAINTEXT://:19092,PLAINTEXT_HOST://:9092'
      KAFKA_ADVERTISED_LISTENERS: 'PLAINTEXT://broker-2:19092,PLAINTEXT_HOST://localhost:39092'
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@controller-1:9093,2@controller-2:9093,3@controller-3:9093
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
    depends_on:
      - controller-1
      - controller-2
      - controller-3

  broker-3:
    image: apache/kafka:latest
    container_name: broker-3
    ports:
      - 49092:9092
    environment:
      KAFKA_NODE_ID: 6
      KAFKA_PROCESS_ROLES: broker
      KAFKA_LISTENERS: 'PLAINTEXT://:19092,PLAINTEXT_HOST://:9092'
      KAFKA_ADVERTISED_LISTENERS: 'PLAINTEXT://broker-3:19092,PLAINTEXT_HOST://localhost:49092'
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@controller-1:9093,2@controller-2:9093,3@controller-3:9093
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
    depends_on:
      - controller-1
      - controller-2
      - controller-3
```

### bitnami kafka单实例 基于zookeeper

运行1个broker实例, 客户端访问: `bootstrap.servers=localhost:29092`

```
services:
  zookeeper:
    container_name: zookeeper
    image: bitnami/zookeeper:3.8
    ports:
      - "2181:2181"
    environment:
      ALLOW_ANONYMOUS_LOGIN: yes

  kafka:
    container_name: kafka
    image: bitnami/kafka:3.9.0
    ports:
      - "29092:29092"
    environment:
      KAFKA_CFG_NODE_ID: 0
      KAFKA_CFG_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,EXTERNAL://:29092
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,EXTERNAL://localhost:29092
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
    depends_on:
      - zookeeper
```

### bitnami kafka多实例 基于zookeeper

运行3个broker实例, 客户端访问: `bootstrap.servers=localhost:29092,localhost:39092,localhost:49092`

```
services:
  zookeeper:
    container_name: zookeeper
    image: bitnami/zookeeper:3.8
    environment:
      ALLOW_ANONYMOUS_LOGIN: yes
    ports:
      - 2181:2181

  kafka-0:
    container_name: kafka-0
    image: bitnami/kafka:3.9.0
    ports:
      - 29092:29092
    environment:
      KAFKA_CFG_NODE_ID: 0
      KAFKA_CFG_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,EXTERNAL://:29092 #2个listener,一个内部用,一个外部用
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,EXTERNAL:PLAINTEXT #要和listener的定义对应,EXTERNAL是自定义的名称
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka-0:9092,EXTERNAL://localhost:29092 #这个配置告诉其他broker和客户端如何访问这个broker, 要和listener的定义对应
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT #在多个listener的情况下最好明确指定一下用哪个
    depends_on:
      - zookeeper

  kafka-1:
    container_name: kafka-1
    image: bitnami/kafka:3.9.0
    ports:
      - 39092:39092
    environment:
      KAFKA_CFG_NODE_ID: 1
      KAFKA_CFG_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,EXTERNAL://:39092
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka-1:9092,EXTERNAL://localhost:39092
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
    depends_on:
      - zookeeper

  kafka-2:
    container_name: kafka-2
    image: bitnami/kafka:3.9.0
    ports:
      - 49092:49092
    environment:
      KAFKA_CFG_NODE_ID: 2
      KAFKA_CFG_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,EXTERNAL://:49092
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka-2:9092,EXTERNAL://localhost:49092
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
    depends_on:
      - zookeeper
```

### bitnami kafka单实例 基于kraft

运行1个broker实例, 客户端访问: `bootstrap.servers=localhost:29092`

```
services:
  kafka:
    container_name: kafka
    image: 'bitnami/kafka:latest'
    ports:
      - '29092:29092'
    environment:
      - KAFKA_CFG_NODE_ID=0
      - KAFKA_CFG_PROCESS_ROLES=controller,broker
      - KAFKA_CFG_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:29092
      - KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092,EXTERNAL://localhost:29092
      - KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT
      - KAFKA_CFG_CONTROLLER_QUORUM_VOTERS=0@kafka:9093
      - KAFKA_CFG_CONTROLLER_LISTENER_NAMES=CONTROLLER
    volumes:
      - kafka-data:/bitnami/kafka
volumes:
  kafka-data:
```

### bitnami kafka多实例 基于kraft

运行3个broker实例, 客户端访问: `bootstrap.servers=localhost:29092,localhost:39092,localhost:49092`
和apache kafka不同, 没有`KAFKA_KRAFT_CLUSTER_ID`这行会启动失败.
如果还有问题可以把 `ALLOW_PLAINTEXT_LISTENER: yes` 注释放开试试, 一开始没加这条不成功, 成功之后去掉也没复现一开始的失败, 可能不是必须的吧

```
services:
  kafka-0:
    container_name: kafka-0
    image: bitnami/kafka:3.9.0
    ports:
      - 29092:29092
    environment:
      KAFKA_CFG_NODE_ID: 0
      KAFKA_CFG_PROCESS_ROLES: broker,controller
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:29092
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka-0:9092,CONTROLLER://kafka-0:9093,EXTERNAL://localhost:29092
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_KRAFT_CLUSTER_ID: kafka-cluster
      # ALLOW_PLAINTEXT_LISTENER: yes

  kafka-1:
    container_name: kafka-1
    image: bitnami/kafka:3.9.0
    ports:
      - 39092:39092
    environment:
      KAFKA_CFG_NODE_ID: 1
      KAFKA_CFG_PROCESS_ROLES: broker,controller
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:39092
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka-1:9092,CONTROLLER://kafka-1:9093,EXTERNAL://localhost:39092
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_KRAFT_CLUSTER_ID: kafka-cluster
      # ALLOW_PLAINTEXT_LISTENER: yes

  kafka-2:
    container_name: kafka-2
    image: bitnami/kafka:3.9.0
    ports:
      - 49092:49092
    environment:
      KAFKA_CFG_NODE_ID: 2
      KAFKA_CFG_PROCESS_ROLES: broker,controller
      KAFKA_CFG_CONTROLLER_QUORUM_VOTERS: 0@kafka-0:9093,1@kafka-1:9093,2@kafka-2:9093
      KAFKA_CFG_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_CFG_LISTENERS: PLAINTEXT://:9092,CONTROLLER://:9093,EXTERNAL://:49092
      KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,CONTROLLER:PLAINTEXT,EXTERNAL:PLAINTEXT
      KAFKA_CFG_ADVERTISED_LISTENERS: PLAINTEXT://kafka-2:9092,CONTROLLER://kafka-2:9093,EXTERNAL://localhost:49092
      KAFKA_CFG_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_KRAFT_CLUSTER_ID: kafka-cluster
      # ALLOW_PLAINTEXT_LISTENER: yes
```

### kafka ui

一个 kafka 管理界面

```
services:
  kafka-ui:
    image: provectuslabs/kafka-ui
    container_name: kafka-ui
    ports:
      - "8080:8080"
    environment:
      DYNAMIC_CONFIG_ENABLED: "true"
```

或者 

```
docker run -it -d --network docker_default -p 8080:8080 -e DYNAMIC_CONFIG_ENABLED=true provectuslabs/kafka-ui
```

## 消息传递语义

- 最多一次（At Most Once）
  - 每条消息最多被处理一次，**可能丢失，但绝不会重复**。
  - 发送消息后不等待确认，不做重试。
  - **性能最高，可靠性最低**。
  - 适用场景：日志采集、监控数据等可容忍丢失的业务。

- 最少一次（At Least Once）
  - 每条消息至少被处理一次，**不会丢失，但可能重复**。
  - 生产端或消费者失败会触发重试。
  - **需要消费者具备幂等性处理能力**（如去重逻辑）。
  - 适用场景：数据采集、任务处理等对完整性要求较高但可容忍重复的业务。

- 精确一次（Exactly Once）
  - 每条消息被处理 **且仅处理一次**，**既不重复也不丢失**。
  - Kafka 提供端到端事务机制实现。
  - **性能开销大，实现复杂**。
  - 适用场景：订单系统、金融交易、库存更新等对一致性要求极高的系统。

Kafka 中实现不同语义的配置方式：

| 配置/行为                        | 最多一次 (At Most Once) | 最少一次 (At Least Once) | 精确一次 (Exactly Once)          |
| -------------------------------- | ----------------------- | ------------------------ | -------------------------------- |
| `acks=0`                         | ✅                       | ❌                        | ❌                                |
| `acks=all`                       | ❌                       | ✅                        | ✅                                |
| `enable.idempotence=true`        | ❌                       | ✅（建议）                | ✅（必须）                        |
| `transactional.id`               | ❌                       | ❌                        | ✅（必须）                        |
| 使用事务 API (`Begin/Commit`)    | ❌                       | ❌                        | ✅                                |
| 生产端重试机制                   | ❌                       | ✅                        | ✅                                |
| 消费端幂等处理                   | ❌                       | ✅（建议）                | ✅（必须）                        |
| `isolation.level=read_committed` | ❌                       | 可选                     | ✅（必须）                        |
| `enable.auto.commit=false`       | ❌                       | ✅（建议）                | ✅（必须，事务中手动提交 offset） |

## 验证

Note:

1.测试前确保没有脏数据, 删除一下sqlite的db文件, 或者重启一下docker容器, 容器没有做数据持久化, 重启便可清除. 

2.处理每条消息加了1秒的延迟

### test-exactly-once

- [x] 精确一次
- [x] 消费9条消息需要9秒
- [x] 优雅退出

```
cd test-exactly-once
go run cmd/consumer/main.go
go run cmd/producer/main.go
```

### test-consumer-parallel

验证消费者并行的处理消息

- [x] 精确一次
- [x] 消费9条数据需要3秒
- [x] 优雅退出

```
cd test-exactly-once
go run cmd/consumer/main.go
go run cmd/producer/main.go
```

### test-consumer-concurrent

验证3费者并行的处理消息的同时每个消费者内部并发的处理消息

- [x] 精确一次
- [x] 消费9条数据需要1秒
- [x] 优雅退出

```
cd test-consumer-concurrent
go run cmd/consumer/main.go
go run cmd/producer/main.go
```
