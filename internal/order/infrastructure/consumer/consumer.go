package consumer

import (
	"context"
	"encoding/json"

	"github.com/QMEOWQ/go-microservice-proj/common/broker"
	"github.com/QMEOWQ/go-microservice-proj/order/app"
	"github.com/QMEOWQ/go-microservice-proj/order/app/command"
	domain "github.com/QMEOWQ/go-microservice-proj/order/domain/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreated, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("fail to consume: queue=%s, err=%v", q.Name, err)
	}

	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(msg)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(msg amqp.Delivery) {
	o := &domain.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		if err := json.Unmarshal(msg.Body, o); err != nil {
			logrus.Infof("failed to unmarshall msg to order, err=%v", err)
			_ = msg.Nack(false, false)
			return
		}
	}

	_, err := c.app.Commands.UpdateOrder.Handle(context.Background(), command.UpdateOrder{
		Order: o,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			if err := order.IsPaid(); err != nil {
				return nil, err
			}
			return order, nil
		},
	})
	if err != nil {
		logrus.Infof("error updating order, orderID = %s, err = %v", o.ID, err)
		// TODO: retry
		return
	}

	_ = msg.Ack(false)
	logrus.Info("order consume paid event success!")
}
