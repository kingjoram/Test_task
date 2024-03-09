package usecase

import (
	"errors"
	"log/slog"
	"test/configs"
	"test/repository/postgres"
)

var ErrUncorrectInput = errors.New("got uncorrect input string")

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
	if len(short) != 10 {
		return "", ErrUncorrectInput
	}
	return "", nil
}
