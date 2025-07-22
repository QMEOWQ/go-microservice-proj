package ports

import (
	context "context"

	"github.com/QMEOWQ/go-microservice-proj/common/genproto/orderpb"
	"github.com/QMEOWQ/go-microservice-proj/order/app"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app : app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	//TODO 
	panic("not implemented yet")
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	//TODO 
	panic("not implemented yet")
}

func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*emptypb.Empty, error) {
	//TODO 
	panic("not implemented yet")
}
