package main

import (
	"flag"
	"log/slog"
	"os"
	"test/configs"
	"test/delivery"
	"test/repository/postgres"
	"test/usecase"

	"github.com/joho/godotenv"
)

func main() {
	var path string
	flag.StringVar(&path, "log_path", "log.log", "Путь к логу")
	logFile, _ := os.Create(path)
	lg := slog.New(slog.NewJSONHandler(logFile, nil))
	lg.Info("start main")

	err := godotenv.Load()
	if err != nil {
		lg.Error("no .env file found")
	}
	lg.Info("load .env")

	config, err := configs.ReadConfig()
	if err != nil {
		lg.Error("read config error", "err", err.Error())
		return
	}

	var db postgres.IDbRepo
	switch config.Db {
	case "postgres":
		db, err = postgres.GetPostgreRepo(config, lg)
	}
	if err != nil {
		lg.Error("cant create repo")
		return
	}

	core := usecase.GetCore(config, lg, db)
	api := delivery.GetApi(core, lg, config)

	api.ListenAndServe()
}
