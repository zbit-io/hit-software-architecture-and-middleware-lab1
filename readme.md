# 实验1：简易消息中间件开发
1. 不可借助开源消息中间件框架
2. 掌握事件驱动架构风格的原理
3. 使用观察者/被观察者模式或发布/订阅模式进行设计和实现
4. 模拟实现一种简单的消息中间件，至少能达到单机系统的功能解耦目的
5. 分析消息中间件的吞吐率等非功能指标
6. 应结合可提供网络服务的管理系统(或其他实际软件系统)进行分析，给出哪些场景可以运用该消息中间件
应该出设计过程和实现细节


## 架构设计
``` lua
+-------------------+
|      Publisher    |
+-------------------+
         |
         | publish
         v
+-------------------+
|       Broker      |
|-------------------|
|                   |
|  +-------------+  |
|  | Subscriber  |  |
|  +-------------+  |
|                   |
+-------------------+
         |
         | notify
         v
+-------------------+
|    Subscriber 1   |
+-------------------+
         |
         | notify
         v
+-------------------+
|    Subscriber 2   |
+-------------------+

```

```
- /cmd/main.go：项目的入口文件，启动 WebSocket 服务器并初始化 Broker 和 Publisher。
/pkg/broker/broker.go：实现消息 Broker 的逻辑，包括管理连接、发布和分发消息的功能。
/pkg/publisher/publisher.go：实现 Publisher 的功能，负责获取用户消息并将其发送到 Broker。
/pkg/subscriber/subscriber.go：实现 Subscriber 的功能，负责接收 Broker 转发的消息并处理。
/pkg/message/message.go：定义消息的结构体和相关的逻辑，例如序列化和反序列化。
/internal/websocket/websocket.go：专门处理 WebSocket 连接的逻辑，包括连接管理和消息传输。
/configs/config.yaml：项目配置文件，用于存储环境变量和其他配置信息。
/scripts/setup.sh：项目的初始化脚本，帮助设置开发环境或部署环境。
/docs：文档目录，包含架构设计、API 文档等。
```