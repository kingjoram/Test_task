package usecase

import (
	"log/slog"
	"test/configs"
	"test/repository/postgres"
)

type ICore interface {
	GetShort(long string) (string, error)
	GetLong(short string) (string, error)
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

func (core *Core) GetShort(long string) (string, error) {

	return "", nil
}

func (core *Core) GetLong(short string) (string, error) {

	return "", nil
}
