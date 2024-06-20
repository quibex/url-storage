package app

import (
	"log/slog"
	grpcapp "url-storage/internal/app/grpc"
	"url-storage/internal/cache/redis"
	"url-storage/internal/config"
	"url-storage/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	const op = "app.New"

	log = log.With(slog.String("op", op))

	log.Info("initializing storage")
	storage, err := postgres.New(&cfg.Postgres)
	if err != nil {
		panic(err)
	}
	log.Info("storage initialized")

	log.Info("initializing cache")
	cache, err := redis.New(&cfg.Redis, storage)
	if err != nil {
		panic(err)
	}
	log.Info("cache initialized")

	log.Info("initializing grpc server")
	grpcServer := grpcapp.New(log, cfg.GRPC.Port, cache)
	log.Info("grpc server initialized")

	return &App{
		GRPCSrv: grpcServer,
	}
}
