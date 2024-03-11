package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"test/configs"
	"test/pkg/models"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type RepoPostgre struct {
	db *sql.DB
}

func GetPostgreRepo(config *configs.Config, lg *slog.Logger) (*RepoPostgre, error) {
	lg.Info("creating postgres repo")

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

func (repo *RepoPostgre) InsertUrl(url models.Url) error {
	_, err := repo.db.Exec(
		"INSERT INTO url(long, short) VALUES ($1, $2)", url.Long, url.Short)
	if err != nil {
		return fmt.Errorf("insert url error: %w", err)
	}

	return nil
}

func (repo *RepoPostgre) GetId() (uint64, error) {
	var id uint64
	err := repo.db.QueryRow("SELECT COUNT(short) + 1 as max FROM url").Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("get id error: %w", err)
	}

	return id, nil
}

func (repo *RepoPostgre) GetShort(long string) (string, error) {
	var short string
	err := repo.db.QueryRow("SELECT short FROM url "+
		"WHERE long = $1", long).Scan(&short)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("get short error: %w", err)
	}

	return short, nil
}

func (repo *RepoPostgre) GetLong(short string) (string, error) {
	var long string
	err := repo.db.QueryRow("SELECT long FROM url "+
		"WHERE short = $1", short).Scan(&long)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}

		return "", fmt.Errorf("get long error: %w", err)
	}

	return short, nil
}
