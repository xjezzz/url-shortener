package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/posgtres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Data struct {
	Id    int
	Alias string
	Url   string
}

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debugger enabled")

	var data Data

	db, err := posgtres.New(context.Background(), cfg.Settings)
	if err != nil {
		fmt.Println(123)
		return
	}

	// Вместо использования QueryRow, используйте QueryRowContext
	err = db.QueryRow(context.Background(), "SELECT id, alias, url FROM url LIMIT 1").Scan(&data.Id, &data.Alias, &data.Url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)

	// TODO: init router: chi

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
