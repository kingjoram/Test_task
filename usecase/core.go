package usecase

import (
	"log/slog"
	"test/configs"
	"test/repository/postgres"
)

type ICore interface {
}

type Core struct {
	lg *slog.Logger
	db postgres.IDbRepo
}

func GetCore(cfg *configs.Config, lg *slog.Logger, db postgres.IDbRepo) *Core {
	lg.Info("creating core")
	return &Core{
		lg: lg,
		db: db,
	}
}
