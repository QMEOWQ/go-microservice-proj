package main

import (
	"log"

	"github.com/QMEOWQ/go-microservice-proj/common/config"
	"github.com/QMEOWQ/go-microservice-proj/common/genproto/orderpb"
	"github.com/QMEOWQ/go-microservice-proj/common/server"
	"github.com/QMEOWQ/go-microservice-proj/ports"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	serviceName := viper.GetString("order.service-name")
	
	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer()
		orderpb.RegisterOrderServiceServer(server, svc)
	})
	
	server.RunHTTPServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, HTTPServer{}, ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

	// 2025/04
	// log.Println("Listenting on port 8085")

	// mux := http.NewServeMux()

	// mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("%v", r.RequestURI)
	// 	_, _ = io.WriteString(w, "<h1>Welcome to Go Miserver</h1>")
	// })

	// mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("%v", r.RequestURI)
	// 	_, _ = io.WriteString(w, "pong")
	// })

	// if err := http.ListenAndServe(":8085", mux); err != nil {
	// 	log.Fatal(err)
	// }
}
