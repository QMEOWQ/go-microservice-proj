package main

import (
	"fmt"
	"net/http"

	"github.com/QMEOWQ/go-microservice-proj/common"
	client "github.com/QMEOWQ/go-microservice-proj/common/client/order"
	"github.com/QMEOWQ/go-microservice-proj/common/tracing"
	"github.com/QMEOWQ/go-microservice-proj/order/app"
	"github.com/QMEOWQ/go-microservice-proj/order/app/command"
	"github.com/QMEOWQ/go-microservice-proj/order/app/query"
	"github.com/QMEOWQ/go-microservice-proj/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app app.Application
	common.BaseResponse
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	_, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
	defer span.End()

	var (
		req  client.CreateOrderRequest
		err  error
		resp struct {
			CustomerID  string `json:"customer_id"`
			OrderID     string `json:"order_id"`
			RedirectURL string `json:"redirect_url"`
		}
	)
	defer func() {
		H.Response(c, err, &resp)
	}()

	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CreateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp.CustomerID = req.CustomerId
	resp.RedirectURL = fmt.Sprintf("http://localhost:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID)
	resp.OrderID = r.OrderID
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	_, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
	defer span.End()

	var (
		err  error
		resp interface{}
	)
	defer func() {
		H.Response(c, err, resp)
	}()

	o, err := H.app.Queries.GetCustomerOrder.Handle(c, query.GetCustomerOrder{
		// OrderID:    "fake-ID",
		// CustomerID: "fake-customer-id",
		OrderID:    orderID,
		CustomerID: customerID,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}
	resp = convertor.NewOrderConvertor().EntityToClient(o)
}
