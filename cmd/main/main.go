package main

import (
	"flag"
	"log/slog"
	"os"
	"test/configs"
	"test/delivery"
	"test/repository"
	"test/repository/memory"
	"test/repository/postgres"
	"test/usecase"
)

func main() {
	var path string
	flag.StringVar(&path, "log_path", "log.log", "Путь к логу")
	logFile, _ := os.Create(path)
	lg := slog.New(slog.NewJSONHandler(logFile, nil))
	lg.Info("start main")

	config, err := configs.ReadConfig()
	if err != nil {
		lg.Error("read config error", "err", err.Error())
		return
	}

	var db repository.IDbRepo
	switch config.Db {
	case "postgres":
		db, err = postgres.GetPostgreRepo(config, lg)
	case "memory":
		db, err = memory.GetMemoryRepo(lg)
	}
	if err != nil {
		lg.Error("cant create repo")
		return
	}
	lg.Info("repo is created")

	core := usecase.GetCore(config, lg, db)
	api := delivery.GetApi(core, lg, config)

	api.ListenAndServe()
}
