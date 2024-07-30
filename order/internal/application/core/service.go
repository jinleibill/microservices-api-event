package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinleibill/web-toolkit-go/egress"
	"github.com/jinleibill/web-toolkit-go/grpc"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"order/config"
	"order/internal/adapters/base"
	grpcServer "order/internal/adapters/grpc"
	"time"
)

type Service struct {
	appFn          func(*Service) error
	Conn           base.Client
	AggregateStore base.Store
	GrpcServer     grpc.Server
}

func NewService(appFn func(*Service) error) *Service {
	return &Service{appFn: appFn}
}

func (s *Service) Run() error {
	dsn := config.GetDataSourceURL()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	s.Conn = base.NewSessionClient(db)
	s.AggregateStore = base.NewSnapshotStore(s.Conn)(base.NewEventStore(s.Conn))

	s.GrpcServer = grpc.NewServer(
		fmt.Sprintf(":%d", config.GetApplicationPort()),
		grpc.WithUnaryServerInterceptors(grpcServer.SessionUnaryInterceptor(s.Conn)),
	)

	err = s.appFn(s)
	if err != nil {
		return err
	}

	waiter := egress.NewWaiter()

	waiter.Add(s.waitForGrpcServer)

	return waiter.Wait()
}

func (s *Service) waitForGrpcServer(ctx context.Context) (err error) {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		err = s.GrpcServer.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		tCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err = s.GrpcServer.Shutdown(tCtx); err != nil {
			return err
		}

		return nil
	})

	return group.Wait()
}
