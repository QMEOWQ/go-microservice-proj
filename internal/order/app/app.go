package app

import (
	"github.com/QMEOWQ/go-microservice-proj/order/app/command"
	"github.com/QMEOWQ/go-microservice-proj/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	UpdateOrder command.UpdateOrderHandler
}

type Queries struct {
	GetCustomerOrder query.GetCustomerOrderHandler
}
