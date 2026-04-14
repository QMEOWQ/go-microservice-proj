# go_miserver3 项目代码梳理文档

> 架构图、数据流图、函数调用图见 [DIAGRAMS.md](./DIAGRAMS.md)

## 1. 项目概览

微服务架构，包含 Order、Stock、Payment 三个服务，采用 CQRS、端口与适配器模式。

| 服务 | 入口 | 协议 | 端口 |
|------|------|------|------|
| Order | `internal/order/main.go` | HTTP + gRPC | HTTP 8282, gRPC 5002 |
| Stock | `internal/stock/main.go` | gRPC | 5003 |
| Payment | `internal/payment/main.go` | HTTP | 8284 |

---

## 2. 全接口清单

### 2.1 领域层接口

| 包 | 接口 | 方法 |
|----|------|------|
| `order/domain/order` | `Repository` | `Create(ctx, *Order) (*Order, error)`<br>`Get(ctx, id, customerID string) (*Order, error)`<br>`Update(ctx, *Order, updateFn) error` |
| `stock/domain/stock` | `Repository` | `GetItems(ctx, ids []string) ([]*orderpb.Item, error)` |

### 2.2 应用层接口

| 包 | 接口 | 说明 |
|----|------|------|
| `order/app/query` | `StockService` | `CheckIfItemsInStock(ctx, items) (*stockpb.CheckIfItemsInStockResponse, error)`<br>`GetItems(ctx, itemIDs) ([]*orderpb.Item, error)` |
| `common/decorator` | `CommandHandler[C, R]` | `Handle(ctx, cmd C) (R, error)` |
| `common/decorator` | `QueryHandler[Q, R]` | `Handle(ctx, query Q) (R, error)` |
| `common/decorator` | `MetricsClient` | `Inc(key string, value int)` |
| `common/discovery` | `Registry` | `Register(ctx, instanceID, serviceName, hostPort)`<br>`Deregister(ctx, instanceID, serviceName)`<br>`Discover(ctx, serviceName)`<br>`HealthCheck(instanceID, serviceName)` |
| `order/ports` | `ServerInterface` | `PostCustomerCustomerIDOrders(c, customerID)`<br>`GetCustomerCustomerIDOrdersOrderID(c, customerID, orderID)` |

### 2.3 gRPC 服务接口（proto 生成）

| 服务 | RPC |
|------|-----|
| `OrderService` | `CreateOrder`, `GetOrder`, `UpdateOrder` |
| `StockService` | `GetItems`, `CheckIfItemsInStock` |

---

## 3. 全函数清单

### 3.1 Order 服务

| 文件 | 函数 | 说明 |
|------|------|------|
| `main.go` | `init()` | 初始化日志、Viper 配置 |
| `main.go` | `main()` | 启动 Order 服务（gRPC + HTTP） |
| `http.go` | `PostCustomerCustomerIDOrders` | HTTP POST 创建订单 |
| `http.go` | `GetCustomerCustomerIDOrdersOrderID` | HTTP GET 查询订单 |
| `service/application.go` | `NewApplication` | 组装应用（含 Stock gRPC 客户端） |
| `service/application.go` | `newApplication` | 内部应用组装 |
| `app/app.go` | - | 定义 `Application`, `Commands`, `Queries` 结构体 |
| `app/command/create_order.go` | `NewCreateOrderHandler` | 创建订单 Handler 工厂 |
| `app/command/create_order.go` | `Handle` | 执行创建订单 |
| `app/command/create_order.go` | `validate` | 校验并调用 Stock 检查库存 |
| `app/command/create_order.go` | `packItems` | 合并同 ID 商品数量 |
| `app/command/update_order.go` | `NewUpdateOrderHandler` | 更新订单 Handler 工厂 |
| `app/command/update_order.go` | `Handle` | 执行更新订单 |
| `app/query/get_customer_order.go` | `NewGetCustomerOrderHandler` | 查询订单 Handler 工厂 |
| `app/query/get_customer_order.go` | `Handle` | 执行查询订单 |
| `domain/order/repository.go` | `NotFoundError.Error` | 错误信息 |
| `domain/order/order.go` | - | 定义 `Order` 结构体 |
| `adapters/order_inmem_repository.go` | `NewMemoryOrderRepository` | 内存仓储工厂 |
| `adapters/order_inmem_repository.go` | `Create` | 内存创建订单 |
| `adapters/order_inmem_repository.go` | `Get` | 内存查询订单 |
| `adapters/order_inmem_repository.go` | `Update` | 内存更新订单 |
| `adapters/grpc/stock_grpc.go` | `NewStockGRPC` | Stock gRPC 适配器工厂 |
| `adapters/grpc/stock_grpc.go` | `CheckIfItemsInStock` | 调用 Stock 检查库存 |
| `adapters/grpc/stock_grpc.go` | `GetItems` | 调用 Stock 获取商品 |
| `ports/grpc.go` | `NewGRPCServer` | gRPC 服务工厂 |
| `ports/grpc.go` | `CreateOrder` | gRPC CreateOrder（未实现，panic） |
| `ports/grpc.go` | `GetOrder` | gRPC GetOrder（未实现，panic） |
| `ports/grpc.go` | `UpdateOrder` | gRPC UpdateOrder（未实现，panic） |

### 3.2 Stock 服务

| 文件 | 函数 | 说明 |
|------|------|------|
| `main.go` | `init()` | 初始化日志、Viper 配置 |
| `main.go` | `main()` | 启动 Stock gRPC 服务 |
| `service/application.go` | `NewApplication` | 组装应用 |
| `app/app.go` | - | 定义 `Application`, `Commands`, `Queries` |
| `app/query/get_items.go` | `NewGetItemsHandler` | GetItems Handler 工厂 |
| `app/query/get_items.go` | `Handle` | 执行获取商品 |
| `app/query/check_if_items_in_stock.go` | `NewCheckIfItemsInStockHandler` | 检查库存 Handler 工厂 |
| `app/query/check_if_items_in_stock.go` | `Handle` | 执行检查库存 |
| `domain/stock/repository.go` | `NotFoundError.Error` | 错误信息 |
| `adapters/stock_inmem_repository.go` | `NewMemoryStockRepository` | 内存仓储工厂 |
| `adapters/stock_inmem_repository.go` | `GetItems` | 内存获取商品 |
| `ports/grpc.go` | `NewGRPCServer` | gRPC 服务工厂 |
| `ports/grpc.go` | `GetItems` | gRPC GetItems |
| `ports/grpc.go` | `CheckIfItemsInStock` | gRPC CheckIfItemsInStock |

### 3.3 Payment 服务

| 文件 | 函数 | 说明 |
|------|------|------|
| `main.go` | `init()` | 初始化日志、Viper 配置 |
| `main.go` | `main()` | 连接 RabbitMQ，启动 HTTP 服务 |
| `http.go` | `NewPaymentHandler` | Payment Handler 工厂 |
| `http.go` | `RegisterRoutes` | 注册路由 |
| `http.go` | `handleWebhook` | Stripe webhook 处理（仅打日志） |

### 3.4 Common 公共库

| 文件 | 函数 | 说明 |
|------|------|------|
| `config/viper.go` | `NewViperConfig` | 加载 global.yaml |
| `server/grpc.go` | `init()` | 替换 gRPC 日志 |
| `server/grpc.go` | `RunGRPCServer` | 按服务名启动 gRPC |
| `server/grpc.go` | `RunGRPCServerOnAddr` | 按地址启动 gRPC |
| `server/http.go` | `RunHTTPServer` | 按服务名启动 HTTP |
| `server/http.go` | `RunHTTPServerOnAddr` | 按地址启动 HTTP |
| `client/grpc.go` | `NewStockGRPCClient` | 创建 Stock gRPC 客户端 |
| `client/grpc.go` | `grpcDialOpts` | gRPC 拨号选项 |
| `discovery/discovery.go` | `GenerateInstanceID` | 生成实例 ID |
| `discovery/grpc.go` | `RegisterToConsul` | 注册到 Consul |
| `discovery/consul/consul.go` | `New` | 创建 Consul Registry |
| `discovery/consul/consul.go` | `Register` | 注册服务 |
| `discovery/consul/consul.go` | `Deregister` | 注销服务 |
| `discovery/consul/consul.go` | `Discover` | 服务发现 |
| `discovery/consul/consul.go` | `HealthCheck` | TTL 健康检查 |
| `broker/rabbitmq.go` | `Connect` | 连接 RabbitMQ，声明 Exchange |
| `broker/event.go` | - | 定义 `EventOrderCreated`, `EventOrderPaid` |
| `logging/logrus.go` | `Init` | 初始化 Logrus |
| `logging/logrus.go` | `SetFormatter` | 设置日志格式 |
| `metrics/todo_metrics.go` | `TodoMetrics.Inc` | 占位实现 |
| `decorator/command.go` | `ApplyCommandDecorators` | 应用 Command 装饰器 |
| `decorator/query.go` | `ApplyQueryDecorators` | 应用 Query 装饰器 |
| `decorator/logging.go` | `queryLoggingDecorator.Handle` | 日志装饰 |
| `decorator/logging.go` | `generateActionName` | 生成动作名 |
| `decorator/metrics.go` | `queryMetricsDecorator.Handle` | 指标装饰 |

---

## 4. 全 API 端点

### 4.1 Order HTTP API（OpenAPI）

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/customer/:customerID/orders/` | 创建订单 |
| GET | `/api/customer/:customerID/orders/:orderID` | 查询订单 |

### 4.2 Order gRPC API

| RPC | 状态 |
|-----|------|
| CreateOrder | 未实现（panic） |
| GetOrder | 未实现（panic） |
| UpdateOrder | 未实现（panic） |

### 4.3 Stock gRPC API

| RPC | 说明 |
|-----|------|
| GetItems | 根据 ItemIDs 返回商品列表 |
| CheckIfItemsInStock | 检查商品库存，返回可用商品 |

### 4.4 Payment HTTP API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/webhook` | Stripe webhook（仅打日志） |

### 4.5 公共端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/ping` | 健康检查，返回 "pong!" |

---

## 5. 启动后变量与运行状态

### 5.1 Order 服务启动流程

```
init()
  ├── logging.Init()           → 设置 Logrus JSON/Text 格式、Debug 级别
  └── config.NewViperConfig()  → 加载 ../common/config/global.yaml

main()
  ├── serviceName = viper["order.service-name"]  → "order"
  ├── ctx, cancel = context.WithCancel(...)
  ├── application, cleanup = service.NewApplication(ctx)
  │     ├── stockClient = NewStockGRPCClient(ctx)  → 连接 stock:5003
  │     ├── stockGRPC = NewStockGRPC(stockClient)
  │     ├── orderRepo = NewMemoryOrderRepository()  → 含 1 条假数据
  │     ├── logger = logrus.NewEntry(...)
  │     ├── metricClient = TodoMetrics{}
  │     └── 返回 Application{Commands, Queries}
  ├── deregisterFunc = RegisterToConsul(ctx, "order")
  │     ├── registry = consul.New("127.0.0.1:8500")
  │     ├── instanceID = "order-<随机数>"
  │     ├── registry.Register(..., "127.0.0.1:5002")
  │     └── goroutine: 每 1s HealthCheck
  ├── goroutine: RunGRPCServer("order", ...)  → 监听 127.0.0.1:5002
  └── RunHTTPServer("order", ...)  → 监听 127.0.0.1:8282（阻塞）
```

**Order 运行时变量：**

| 变量 | 类型 | 说明 |
|------|------|------|
| `application` | `app.Application` | Commands(CreateOrder, UpdateOrder) + Queries(GetCustomerOrder) |
| `orderRepo` | `*MemoryOrderRepository` | 内存 store，初始含 fake-ID 订单 |
| `stockGRPC` | `*StockGRPC` | 连接 stock:5003 的 gRPC 客户端 |
| `deregisterFunc` | `func() error` | 退出时从 Consul 注销 |

### 5.2 Stock 服务启动流程

```
init()
  └── 同 Order

main()
  ├── serviceName = "stock"
  ├── serverType = viper["stock.server-to-run"]  → "grpc"
  ├── application = service.NewApplication(ctx)
  │     ├── stockRepo = NewMemoryStockRepository()  → stub 数据 (item_id, item1, item2, item3)
  │     ├── logger, metricsClient
  │     └── 返回 Application{Queries: GetItems, CheckIfItemsInStock}
  ├── deregisterFunc = RegisterToConsul(ctx, "stock")  → 127.0.0.1:5003
  └── RunGRPCServer("stock", ...)  → 监听 127.0.0.1:5003（阻塞）
```

**Stock 运行时变量：**

| 变量 | 类型 | 说明 |
|------|------|------|
| `application` | `app.Application` | Queries(GetItems, CheckIfItemsInStock) |
| `stockRepo` | `*MemoryStockRepository` | store 含 item_id/item1/item2/item3 |

### 5.3 Payment 服务启动流程

```
init()
  └── 同 Order

main()
  ├── serverType = "http"
  ├── ch, closeCh = broker.Connect(user, password, host, port)
  │     └── 声明 Exchange: order_created(direct), order.paid(fanout)
  ├── paymentHandler = NewPaymentHandler()
  └── RunHTTPServer("payment", paymentHandler.RegisterRoutes)  → 监听 127.0.0.1:8284
```

**Payment 运行时变量：**

| 变量 | 类型 | 说明 |
|------|------|------|
| `ch` | `*amqp.Channel` | RabbitMQ 通道 |
| `closeCh` | `func() error` | 关闭连接 |
| `paymentHandler` | `*PaymentHandler` | 空结构体，RegisterRoutes 注册 /api/webhook |

---

## 6. 系统依赖与运行要求

### 6.1 外部依赖

| 组件 | 地址 | 用途 |
|------|------|------|
| Consul | 127.0.0.1:8500 | 服务注册与发现 |
| RabbitMQ | 127.0.0.1:5672 | Payment 连接，声明 Exchange |
| Stock gRPC | 127.0.0.1:5003 | Order 调用库存 |

### 6.2 启动顺序建议

1. `docker-compose up` 启动 Consul、RabbitMQ
2. 启动 Stock（Order 依赖其 gRPC）
3. 启动 Order
4. 启动 Payment

### 6.3 配置文件 (global.yaml)

| 键 | 值 |
|----|-----|
| order.service-name | order |
| order.http-addr | 127.0.0.1:8282 |
| order.grpc-addr | 127.0.0.1:5002 |
| stock.service-name | stock |
| stock.grpc-addr | 127.0.0.1:5003 |
| payment.service-name | payment |
| payment.http-addr | 127.0.0.1:8284 |
| consul.addr | 127.0.0.1:8500 |
| rabbitmq.* | user/password/host/port |

---

## 7. 数据流示意

```
[HTTP Client] 
    → POST /api/customer/:id/orders/ 
    → HTTPServer.PostCustomerCustomerIDOrders 
    → CreateOrderHandler.Handle 
    → orderRepo.Create + stockGRPC.CheckIfItemsInStock
    → MemoryOrderRepository / Stock gRPC

[Order gRPC Client] 
    → Stock:5003 GetItems / CheckIfItemsInStock 
    → StockGRPCServer 
    → GetItemsHandler / CheckIfItemsInStockHandler 
    → MemoryStockRepository
```

---

## 8. 注意事项

1. **Order gRPC**：CreateOrder/GetOrder/UpdateOrder 当前为 `panic("not implemented yet")`，实际使用 HTTP API。
2. **Payment**：`internal/payment/main.go` 使用 `package payment`，若需 `go run` 需确认包名与构建方式。
3. **Viper 配置路径**：`AddConfigPath("../common/config")` 依赖工作目录，建议从项目根或 `internal/common` 运行。
4. **Stock CheckIfItemsInStock**：未做真实库存扣减，仅返回请求中的商品列表。
