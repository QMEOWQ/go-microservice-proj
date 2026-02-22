package service

import (
	"context"

	grpcClient "github.com/QMEOWQ/go-microservice-proj/common/client"
	"github.com/QMEOWQ/go-microservice-proj/common/metrics"
	"github.com/QMEOWQ/go-microservice-proj/order/adapters"
	"github.com/QMEOWQ/go-microservice-proj/order/adapters/grpc"
	"github.com/QMEOWQ/go-microservice-proj/order/app"
	"github.com/QMEOWQ/go-microservice-proj/order/app/command"
	"github.com/QMEOWQ/go-microservice-proj/order/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGRPC(stockClient)

	return newApplication(ctx, stockGRPC), func() {
		_ = closeStockClient()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}

	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, logger, metricClient),
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
