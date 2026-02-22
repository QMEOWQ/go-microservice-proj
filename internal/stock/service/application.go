package service

import (
	"context"

	"github.com/QMEOWQ/go-microservice-proj/common/metrics"
	"github.com/QMEOWQ/go-microservice-proj/stock/adapters"
	"github.com/QMEOWQ/go-microservice-proj/stock/app"
	"github.com/QMEOWQ/go-microservice-proj/stock/app/query"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	stockRepo := adapters.NewMemoryStockRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
		},
	}
}
