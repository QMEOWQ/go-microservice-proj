go mod init github.com/QMEOWQ/go-microservice-proj

20250622: finish initial http&grpc server

test internal/order :

go run main.go http.go:

```
success msg :
[GIN-debug] Listening and serving HTTP on 127.0.0.1:8282
time="2025-06-22T23:21:31+08:00" level=info msg="Starting gRPC server, Listening 127.0.0.1:5002"
```
