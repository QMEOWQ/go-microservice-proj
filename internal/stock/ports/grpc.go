package ports

import (
	context "context"

	"github.com/QMEOWQ/go-microservice-proj/common/genproto/stockpb"
	"github.com/QMEOWQ/go-microservice-proj/stock/app"
	"github.com/QMEOWQ/go-microservice-proj/stock/app/query"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIDs: request.ItemIDs})
	if err != nil {
		return nil, err
	}
	return &stockpb.GetItemsResponse{Items: items}, nil

	// logrus.Info("rpc_request_in, stock.GetItems")
	// defer func() {
	// 	logrus.Info("rpc_request_out, stock.GetItems")
	// }()

	// fake := []*orderpb.Item{
	// 	{
	// 		ID: "fake-item-from-GetItems",
	// 	},
	// }
	// return &stockpb.GetItemsResponse{Items: fake}, nil
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{Items: request.Items})
	if err != nil {
		return nil, err
	}
	return &stockpb.CheckIfItemsInStockResponse{
		InStock: 1,
		Items:   items,
	}, nil

	// logrus.Info("rpc_request_in, stock.CheckIfItemsInStock")
	// defer func() {
	// 	logrus.Info("rpc_request_out, stock.CheckIfItemsInStock")
	// }()

	// return nil, nil

	// fake := []*orderpb.Item{
	// 	{
	// 		ID: "fake-item-from-CheckIfItemsInStock",
	// 	},
	// }
	// return &stockpb.CheckIfItemsInStockResponse{Items: fake}, nil
}
