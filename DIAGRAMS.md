# go_miserver3 架构与流程图

> 使用 Mermaid 语法，可在 GitHub、VS Code、Typora 等支持 Mermaid 的编辑器中渲染。

---

## 1. 整体系统架构图

```mermaid
flowchart TB
    subgraph External["外部依赖"]
        Consul["Consul\n127.0.0.1:8500"]
        RabbitMQ["RabbitMQ\n127.0.0.1:5672"]
    end

    subgraph Order["Order 服务 :8282 HTTP / :5002 gRPC"]
        OH[HTTP Server]
        OG[gRPC Server]
        OApp[Application]
        OApp --> OCmd[Commands]
        OApp --> OQ[Queries]
        OCmd --> CreateOrder[CreateOrder]
        OCmd --> UpdateOrder[UpdateOrder]
        OQ --> GetOrder[GetCustomerOrder]
    end

    subgraph Stock["Stock 服务 :5003 gRPC"]
        SG[gRPC Server]
        SApp[Application]
        SApp --> SQ[Queries]
        SQ --> GetItems[GetItems]
        SQ --> CheckStock[CheckIfItemsInStock]
    end

    subgraph Payment["Payment 服务 :8284 HTTP"]
        PH[HTTP Server]
        Webhook[handleWebhook]
    end

    subgraph Common["Common 公共库"]
        Config[config/viper]
        Discovery[discovery/consul]
        Broker[broker/rabbitmq]
        Decorator[decorator]
    end

    Client1[HTTP Client] --> OH
    Client2[gRPC Client] --> OG
    Stripe[Stripe Webhook] --> PH

    Order -->|gRPC| Stock
    Order --> Consul
    Stock --> Consul
    Payment --> RabbitMQ
    Order --> Config
    Stock --> Config
    Payment --> Config
```

---

## 2. 数据流图

### 2.1 创建订单数据流

```mermaid
flowchart LR
    subgraph Client["客户端"]
        Req["POST /api/customer/:id/orders/\n{ customerID, items }"]
    end

    subgraph Order["Order 服务"]
        HTTP[HTTPServer.PostCustomerCustomerIDOrders]
        Cmd[CreateOrderHandler.Handle]
        Validate[validate + packItems]
        StockCall[StockGRPC.CheckIfItemsInStock]
        RepoCreate[MemoryOrderRepository.Create]
        Resp["JSON { order_id }"]
    end

    subgraph Stock["Stock 服务"]
        StockRPC[GRPCServer.CheckIfItemsInStock]
        StockHandler[CheckIfItemsInStockHandler]
        StockRepo[MemoryStockRepository.GetItems]
    end

    Req --> HTTP --> Cmd
    Cmd --> Validate
    Validate --> StockCall
    StockCall -->|gRPC| StockRPC
    StockRPC --> StockHandler --> StockRepo
    StockRepo -->|Items| StockCall
    StockCall --> RepoCreate
    RepoCreate --> Resp
```

### 2.2 查询订单数据流

```mermaid
flowchart LR
    subgraph Client["客户端"]
        Req["GET /api/customer/:id/orders/:orderID"]
    end

    subgraph Order["Order 服务"]
        HTTP[HTTPServer.GetCustomerCustomerIDOrdersOrderID]
        Query[GetCustomerOrderHandler.Handle]
        Repo[MemoryOrderRepository.Get]
        Resp["JSON { data: Order }"]
    end

    Req --> HTTP --> Query --> Repo --> Resp
```

### 2.3 Stock 服务数据流（Order 调用）

```mermaid
flowchart TB
    subgraph Order["Order 服务"]
        CreateOrder[CreateOrderHandler]
        StockGRPC[StockGRPC 适配器]
    end

    subgraph Stock["Stock 服务"]
        Port[GRPCServer]
        GetItems[GetItemsHandler]
        CheckStock[CheckIfItemsInStockHandler]
        Repo[MemoryStockRepository]
    end

    CreateOrder -->|CheckIfItemsInStock| StockGRPC
    StockGRPC -->|gRPC :5003| Port
    Port --> CheckStock
    Port --> GetItems
    CheckStock --> Repo
    GetItems --> Repo
    Repo -->|orderpb.Item[]| Port
```

---

## 3. 函数调用图

### 3.1 Order 服务启动与请求调用链

```mermaid
flowchart TB
    subgraph Main["main()"]
        M1[service.NewApplication]
        M2[discovery.RegisterToConsul]
        M3[server.RunGRPCServer]
        M4[server.RunHTTPServer]
    end

    subgraph NewApplication["NewApplication"]
        NA1[grpcClient.NewStockGRPCClient]
        NA2[grpc.NewStockGRPC]
        NA3[newApplication]
    end

    subgraph newApplication["newApplication"]
        nA1[adapters.NewMemoryOrderRepository]
        nA2[command.NewCreateOrderHandler]
        nA3[command.NewUpdateOrderHandler]
        nA4[query.NewGetCustomerOrderHandler]
    end

    subgraph HTTPRequest["HTTP 请求: 创建订单"]
        H1[PostCustomerCustomerIDOrders]
        H2[app.Commands.CreateOrder.Handle]
        H3[createOrderHandler.Handle]
        H4[validate]
        H5[packItems]
        H6[stockGRPC.CheckIfItemsInStock]
        H7[orderRepo.Create]
    end

    M1 --> NA1 --> NA2 --> NA3
    NA3 --> nA1
    NA3 --> nA2
    NA3 --> nA3
    NA3 --> nA4

    H1 --> H2 --> H3
    H3 --> H4 --> H5
    H4 --> H6
    H4 --> H7
```

### 3.2 Stock 服务函数调用链

```mermaid
flowchart TB
    subgraph Main["main()"]
        M1[service.NewApplication]
        M2[discovery.RegisterToConsul]
        M3[server.RunGRPCServer]
    end

    subgraph NewApplication["NewApplication"]
        NA1[adapters.NewMemoryStockRepository]
        NA2[query.NewCheckIfItemsInStockHandler]
        NA3[query.NewGetItemsHandler]
    end

    subgraph GRPCRequest["gRPC 请求"]
        G1[GRPCServer.GetItems]
        G2[GRPCServer.CheckIfItemsInStock]
        G3[GetItemsHandler.Handle]
        G4[CheckIfItemsInStockHandler.Handle]
        G5[stockRepo.GetItems]
    end

    M1 --> NA1
    M1 --> NA2
    M1 --> NA3

    G1 --> G3 --> G5
    G2 --> G4 --> G5
```

### 3.3 Payment 服务函数调用链

```mermaid
flowchart TB
    subgraph Main["main()"]
        M1[broker.Connect]
        M2[NewPaymentHandler]
        M3[server.RunHTTPServer]
    end

    subgraph Connect["broker.Connect"]
        C1[amqp.Dial]
        C2[conn.Channel]
        C3[ch.ExchangeDeclare order_created]
        C4[ch.ExchangeDeclare order.paid]
    end

    subgraph HTTPRequest["HTTP 请求"]
        H1[RegisterRoutes]
        H2[handleWebhook]
    end

    M1 --> C1 --> C2 --> C3 --> C4
    M2 --> H1
    H1 -->|POST /api/webhook| H2
```

### 3.4 装饰器调用链（Command/Query）

```mermaid
flowchart LR
    subgraph Decorator["装饰器链"]
        D1[queryLoggingDecorator]
        D2[queryMetricsDecorator]
        D3[createOrderHandler]
    end

    Request[Request] --> D1
    D1 --> D2
    D2 --> D3
    D3 --> Response[Response]
```

---

## 4. 分层架构图（端口与适配器）

```mermaid
flowchart TB
    subgraph Ports["端口 Ports"]
        HTTPPort[HTTP ServerInterface]
        GRPCPort[gRPC Server]
    end

    subgraph Application["应用层 App"]
        Commands[Commands]
        Queries[Queries]
    end

    subgraph Domain["领域层 Domain"]
        OrderDomain[Order 实体]
        StockDomain[Stock 仓储接口]
        OrderRepo[Repository 接口]
    end

    subgraph Adapters["适配器 Adapters"]
        OrderInMem[MemoryOrderRepository]
        StockInMem[MemoryStockRepository]
        StockGRPC[StockGRPC 客户端]
    end

    HTTPPort --> Commands
    HTTPPort --> Queries
    GRPCPort --> Commands
    GRPCPort --> Queries

    Commands --> OrderRepo
    Commands --> StockGRPC
    Queries --> OrderRepo
    Queries --> StockGRPC

    OrderRepo -.->|实现| OrderInMem
    StockGRPC -->|gRPC| StockService
    StockService --> StockDomain
    StockDomain -.->|实现| StockInMem
```

---

## 5. 服务依赖关系图

```mermaid
flowchart LR
    subgraph Services["服务"]
        Order[Order]
        Stock[Stock]
        Payment[Payment]
    end

    subgraph Common["Common"]
        Config[config]
        Discovery[discovery]
        Server[server]
        Client[client]
        Broker[broker]
        Logging[logging]
        Decorator[decorator]
        Metrics[metrics]
    end

    subgraph External["外部"]
        Consul[Consul]
        RabbitMQ[RabbitMQ]
    end

    Order --> Config
    Order --> Discovery
    Order --> Server
    Order --> Client
    Order --> Logging
    Order --> Decorator
    Order --> Metrics
    Order -->|gRPC| Stock

    Stock --> Config
    Stock --> Discovery
    Stock --> Server
    Stock --> Logging
    Stock --> Decorator
    Stock --> Metrics

    Payment --> Config
    Payment --> Server
    Payment --> Broker
    Payment --> Logging

    Discovery --> Consul
    Broker --> RabbitMQ
```
