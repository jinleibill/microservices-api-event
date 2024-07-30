package grpc

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"order/internal/adapters/base"
)

func SessionUnaryInterceptor(client base.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		tx := client.GetTx().WithContext(ctx)

		defer func() {
			p := recover()
			switch {
			case p != nil:
				tx.Rollback()
				if tx.Error != nil {
					log.Printf("error while rolling back the rpc request transaction during panic: %v", tx.Error.Error())
				}
				panic(p)
			case err != nil:
				tx.Rollback()
				if tx.Error != nil {
					log.Printf("error while rolling back the rpc request transaction: %v", tx.Error.Error())
				}
			default:
				tx.Commit()
				if tx.Error != nil {
					log.Printf("error while committing the rpc request transaction: %v", tx.Error.Error())
				}
			}
		}()

		return handler(ctx, req)
	}
}
