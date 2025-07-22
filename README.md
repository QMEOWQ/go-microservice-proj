go mod init github.com/QMEOWQ/go-microservice-proj/dir

replace github.com/QMEOWQ/go-microservice-proj/common => ../common

20250622: finish initial http&grpc server

test internal/order :

go run main.go http.go:

```
success msg :
[GIN-debug] Listening and serving HTTP on 127.0.0.1:8282
time="2025-06-22T23:21:31+08:00" level=info msg="Starting gRPC server, Listening 127.0.0.1:5002"
```

app inject needs pkg:

```
go1.24.3 install github.com/air-verse/air@latest
```

using air pkg

```
air init
air .

msg:
watching .
watching app
watching ports
watching service
!exclude tmp
building...
running...
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /api/customer/:customerID/orders/ --> github.com/QMEOWQ/go-microservice-proj/order/ports.(*ServerInterfaceWrapper).PostCustomerCustomerIDOrders-fm (1 handlers)
[GIN-debug] GET    /api/customer/:customerID/orders/:orderID --> github.com/QMEOWQ/go-microservice-proj/order/ports.(*ServerInterfaceWrapper).GetCustomerCustomerIDOrdersOrderID-fm (1 handlers)
[GIN-debug] GET    /ping                     --> github.com/QMEOWQ/go-microservice-proj/common/server.RunHTTPServerOnAddr.func1 (1 handlers)
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on 127.0.0.1:8282
time="2025-07-22T20:17:30+08:00" level=info msg="Starting gRPC server, Listening 127.0.0.1:5002"
cleaning...
Process Exit with Code: 1
see you again~
```
