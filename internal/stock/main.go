package main

import (
	"github.com/QMEOWQ/go-microservice-proj/common/genproto/stockpb"
	"github.com/QMEOWQ/go-microservice-proj/common/server"
	"github.com/QMEOWQ/go-microservice-proj/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer()
			stockpb.RegisterStockServiceServer(server, svc)
		})

	case "http":
		// 暂时不用

	default:
		panic("unexpected server type")
	}
}
