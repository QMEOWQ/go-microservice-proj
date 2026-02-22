package main

import (
	"context"
	"log"

	"github.com/QMEOWQ/go-microservice-proj/common/config"
	"github.com/QMEOWQ/go-microservice-proj/common/discovery"
	"github.com/QMEOWQ/go-microservice-proj/common/genproto/orderpb"
	"github.com/QMEOWQ/go-microservice-proj/common/server"
	"github.com/QMEOWQ/go-microservice-proj/order/ports"
	"github.com/QMEOWQ/go-microservice-proj/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")

	// log.Panic("I changed this")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application, cleanup := service.NewApplication(ctx)
	defer cleanup() // 在主函数退出时，关闭所有连接

	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{
			app: application,
		}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

	// 2025/04
	// log.Println("Listenting on port 8085")

	// mux := http.NewServeMux()

	// mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("%v", r.RequestURI)
	// 	_, _ = io.WriteString(w, "<h1>Welcome to Go Miserver</h1>")
	// })

	// mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("%v", r.RequestURI)
	// 	_, _ = io.WriteString(w, "pong")
	// })

	// if err := http.ListenAndServe(":8085", mux); err != nil {
	// 	log.Fatal(err)
	// }
}
