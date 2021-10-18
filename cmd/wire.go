//go:build wireinject
// +build wireinject

package cmd

import (
	"context"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/google/wire"
	"spider/configs"
	"spider/internal/scheduler"

	"spider/internal/data"
	"spider/internal/server"
	"spider/internal/service"
)

func InitApp(ctx context.Context, logger *xlog.Logger) *App {
	panic(wire.Build(
		scheduler.ProviderSet,
		server.ProviderSet,
		service.ProviderSet,
		newApp,
		configs.ProviderSet,
		data.ProviderSet,
	))
}
