package main

import (
	"order/internal/adapters/grpc"
	"order/internal/adapters/order"
	"order/internal/application/core"
	"order/internal/application/core/application"
)

func main() {
	err := core.NewService(initService).Run()
	if err != nil {
		panic(err)
	}
}

func initService(s *core.Service) error {
	orderRepoAdapter := order.NewAdapter(s.AggregateStore)

	app := application.NewApplication(orderRepoAdapter)

	grpc.NewAdapter(app, s.Conn).Mount(s.GrpcServer)

	return nil
}
