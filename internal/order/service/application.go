package service

import (
	"context"

	"github.com/QMEOWQ/go-microservice-proj/common/broker"
	grpcClient "github.com/QMEOWQ/go-microservice-proj/common/client"
	"github.com/QMEOWQ/go-microservice-proj/common/metrics"
	"github.com/QMEOWQ/go-microservice-proj/order/adapters"
	"github.com/QMEOWQ/go-microservice-proj/order/adapters/grpc"
	"github.com/QMEOWQ/go-microservice-proj/order/app"
	"github.com/QMEOWQ/go-microservice-proj/order/app/command"
	"github.com/QMEOWQ/go-microservice-proj/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}

	ch, closeCh := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	stockGRPC := grpc.NewStockGRPC(stockClient)
	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeStockClient()
		_ = closeCh()
		_ = ch.Close()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}

	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(
				orderRepo,
				logger,
				metricClient),
		},
	}
}
