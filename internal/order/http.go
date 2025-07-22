package main

import (
	"github.com/QMEOWQ/go-microservice-proj/order/app"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app app.Application
}

func (s HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	panic("implement me")
}

func (s HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	panic("implement me")
}
