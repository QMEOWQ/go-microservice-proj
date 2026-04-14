package app

import "github.com/QMEOWQ/go-microservice-proj/payment/app/command"

type Application struct {
	Commands Commands
}

type Commands struct {
	CreatePayment command.CreatePaymentHandler
}
