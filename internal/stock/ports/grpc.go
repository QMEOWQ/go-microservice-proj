package ports

import (
	context "context"

	"github.com/QMEOWQ/go-microservice-proj/common/genproto/stockpb"
	"github.com/QMEOWQ/go-microservice-proj/stock/app"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	//TODO
	panic("not implemented yet")
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	//TODO
	panic("not implemented yet")
}
