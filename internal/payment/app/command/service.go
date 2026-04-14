package command

import (
	"context"

	"github.com/QMEOWQ/go-microservice-proj/common/genproto/orderpb"
)

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
