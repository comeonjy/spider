package service

import (
	"context"

	"github.com/google/wire"
	"google.golang.org/grpc/metadata"

	"github.com/comeonjy/go-kit/pkg/xlog"
	v1 "spider/api/v1"
	"spider/configs"
	"spider/internal/data"
)

var ProviderSet = wire.NewSet(NewSpiderService)

type SpiderService struct {
	v1.UnimplementedSpiderServer
	conf     configs.Interface
	logger   *xlog.Logger
	taskRepo data.TaskRepo
}

func NewSpiderService(conf configs.Interface, logger *xlog.Logger, taskRepo data.TaskRepo) *SpiderService {
	return &SpiderService{
		conf:     conf,
		taskRepo: taskRepo,
		logger:   logger,
	}
}

func (svc *SpiderService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if mdIn, ok := metadata.FromIncomingContext(ctx); ok {
		mdIn.Get("")
	}
	return ctx, nil
}

func (svc *SpiderService) Ping(ctx context.Context, in *v1.Empty) (*v1.Result, error) {
	return &v1.Result{
		Code:    200,
		Message: "pong",
	}, nil
}