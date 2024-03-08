package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"test/configs"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

type IDbRepo interface {
}

type RepoPostgre struct {
	db *sql.DB
}

func GetPostgreRepo(config *configs.Config, lg *slog.Logger) (*RepoPostgre, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password= %s host=%s port=%d sslmode=%s",
		config.User, config.DbName, config.Password, config.Host, config.Port, config.Sslmode)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		lg.Error("sql open error", "err", err.Error())
		return nil, fmt.Errorf("get db repo: %w", err)
	}
	err = db.Ping()
	if err != nil {
		lg.Error("sql ping error", "err", err.Error())
		return nil, fmt.Errorf("get db repo: %w", err)
	}
	db.SetMaxOpenConns(config.MaxOpenConns)

	postgreDb := RepoPostgre{db: db}

	go postgreDb.pingDb(config.Timer, lg)
	return &postgreDb, nil
}

func (repo *RepoPostgre) pingDb(timer uint32, lg *slog.Logger) {
	for {
		err := repo.db.Ping()
		if err != nil {
			lg.Error("repo db ping error", "err", err.Error())
		}

		time.Sleep(time.Duration(timer) * time.Second)
	}
}
