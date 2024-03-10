package usecase

import (
	"log/slog"
	"test/configs"
	"test/repository"
)

type ICore interface {
}

type Core struct {
	lg *slog.Logger
	db repository.IDbRepo
}

func GetCore(cfg *configs.Config, lg *slog.Logger, db repository.IDbRepo) *Core {
	lg.Info("creating core")
	return &Core{
		lg: lg,
		db: db,
	}
}
