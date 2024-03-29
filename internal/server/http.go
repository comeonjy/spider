package server

import (
	"context"
	"net/http"
	"time"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xmiddleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"spider/api/v1"
	"spider/configs"
	"spider/pkg/consts"
)

func NewHttpServer(ctx context.Context, conf configs.Interface, logger *xlog.Logger) *http.Server {
	mux := runtime.NewServeMux(runtime.WithErrorHandler(xmiddleware.HttpErrorHandler(logger)))
	server := http.Server{
		Addr:              conf.Get().HttpAddr,
		Handler:           xmiddleware.HttpUse(mux, xmiddleware.HttpLogger(consts.TraceName, logger)),
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      2 * time.Second,
	}
	if err := v1.RegisterSpiderHandlerFromEndpoint(ctx, mux, conf.Get().GrpcAddr, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		panic("RegisterSpiderHandlerFromEndpoint" + err.Error())
	}
	return &server
}
