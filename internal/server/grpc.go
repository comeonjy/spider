package server

import (
	"time"

	"github.com/comeonjy/go-kit/pkg/xenv"
	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xmiddleware"
	"github.com/google/wire"
	"google.golang.org/grpc"

	"spider/api/v1"
	"spider/configs"
	"spider/internal/service"
	"spider/pkg/consts"
)

var ProviderSet = wire.NewSet(NewGrpcServer, NewHttpServer)

func NewGrpcServer(srv *service.SpiderService, conf configs.Interface,logger *xlog.Logger) *grpc.Server {
	server := grpc.NewServer(
		grpc.ConnectionTimeout(2*time.Second),
		grpc.ChainUnaryInterceptor(
			xmiddleware.GrpcLogger(consts.TraceName,logger), xmiddleware.GrpcValidate, xmiddleware.GrpcRecover(logger), xmiddleware.GrpcAuth, xmiddleware.GrpcApm(conf.Get().ApmUrl, consts.AppName, consts.AppVersion, xenv.GetEnv(consts.AppEnv))),
	)
	v1.RegisterSpiderServer(server, srv)
	return server
}
