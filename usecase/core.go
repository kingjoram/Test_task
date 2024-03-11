package usecase

import (
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"test/configs"
	"test/pkg/models"
	"test/repository"

	"github.com/sqids/sqids-go"
)

var (
	ErrUncorrectInput = errors.New("got uncorrect input string")
	ErrNotFound       = errors.New("not found")
	AlphabetForShort  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
)

type ICore interface {
	GetShort(long string) (string, error)
	GetLong(short string) (string, error)
}

type Core struct {
	lg    *slog.Logger
	db    repository.IDbRepo
	coder *sqids.Sqids
}

func GetCore(cfg *configs.Config, lg *slog.Logger, db repository.IDbRepo) *Core {
	lg.Info("creating core")
	coder, err := sqids.New(sqids.Options{
		MinLength: 10,
		Alphabet:  AlphabetForShort,
	})
	if err != nil {
		lg.Error("creating core error", "err", err.Error())
		return nil
	}
	return &Core{
		lg:    lg,
		db:    db,
		coder: coder,
	}
}

func (core *Core) GetShort(long string) (string, error) {
	core.lg.Info("start func get short in core")
	match, err := regexp.MatchString(`^[a-z0-9A-z]+://.+`, long)
	if err != nil {
		core.lg.Error("uncorrect regexp", "err", err.Error())
		return "", fmt.Errorf("get short error: %w", err)
	}
	if !match {
		core.lg.Info("uncorrect input string")
		return "", ErrUncorrectInput
	}
	short, err := core.db.GetShort(long)
	if err != nil {
		core.lg.Error("get short encode error", "err", err.Error())
		return "", fmt.Errorf("get short error: %w", err)
	}
	if short != "" {
		return short, nil
	}

	id, err := core.db.GetId()
	if err != nil {
		core.lg.Error("get short get id error", "err", err.Error())
		return "", fmt.Errorf("get short error: %w", err)
	}
	short, err = core.coder.Encode([]uint64{id})
	if err != nil {
		core.lg.Error("get short encode error", "err", err.Error())
		return "", fmt.Errorf("get short error: %w", err)
	}
	short = "http://somedomen.ru/" + short

	err = core.db.InsertUrl(models.Url{Short: short, Long: long})
	if err != nil {
		core.lg.Error("get short save error", "err", err.Error())
		return "", fmt.Errorf("get short error: %w", err)
	}
	return short, nil
}

func (core *Core) GetLong(short string) (string, error) {
	core.lg.Info("start func get long in core")
	match, err := regexp.MatchString(`^http://somedomen.ru/[a-z0-9A-z_]{10}`, short)
	if err != nil {
		core.lg.Error("uncorrect regexp", "err", err.Error())
		return "", fmt.Errorf("get long error: %w", err)
	}
	if !match {
		core.lg.Info("uncorrect input string")
		return "", ErrUncorrectInput
	}

	long, err := core.db.GetLong(short)
	if err != nil {
		core.lg.Error("get long error", "err", err.Error())
		return "", fmt.Errorf("get long error: %w", err)
	}
	if long == "" {
		return "", ErrNotFound
	}

	return long, nil
}
