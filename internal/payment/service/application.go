package service

import (
	"context"

	grpcClient "github.com/QMEOWQ/go-microservice-proj/common/client"
	"github.com/QMEOWQ/go-microservice-proj/common/metrics"
	"github.com/QMEOWQ/go-microservice-proj/payment/adapters"
	"github.com/QMEOWQ/go-microservice-proj/payment/app"
	"github.com/QMEOWQ/go-microservice-proj/payment/app/command"
	"github.com/QMEOWQ/go-microservice-proj/payment/domain"
	"github.com/QMEOWQ/go-microservice-proj/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}

	orderGRPC := adapters.NewOrderGRPC(orderClient)
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}

	// memoryProcessor := processor.NewInmemProcessor()
	// return newApplication(ctx, orderGRPC, memoryProcessor), func() {
	// 	_ = closeOrderClient()
	// }
}

func newApplication(_ context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}
