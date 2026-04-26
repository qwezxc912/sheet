package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/qweq1232/dnd_form/internal/config"
	chardelete "github.com/qweq1232/dnd_form/internal/http_server/handlers/chars/delete"
	"github.com/qweq1232/dnd_form/internal/http_server/handlers/chars/get"
	getall "github.com/qweq1232/dnd_form/internal/http_server/handlers/chars/get_all"
	charsave "github.com/qweq1232/dnd_form/internal/http_server/handlers/chars/save"
	charupdate "github.com/qweq1232/dnd_form/internal/http_server/handlers/chars/update"
	storage "github.com/qweq1232/dnd_form/internal/storage/postgres"
)

const (
	local = "local"
	prod  = "prod"
	dev   = "dev"
)

func main() {
	cfg := config.MustLoad()

	dsn := MustLoadDSN()

	log := setupLogger(cfg.Env)
	log.Info("application is started")

	ctx := context.Background()

	db, err := storage.New(ctx, dsn)
	if err != nil {
		log.Error("failed to connect to databse", slog.Any("err", err))
	}

	log.Info("connected to database")

	r := gin.Default()

	r.POST("/", charsave.New(ctx, log, db))
	r.GET("/:user_id/:id", get.New(ctx, log, db))
	r.GET("/:user_id", getall.New(ctx, log, db))
	r.DELETE("/:id", chardelete.New(ctx, log, db))
	r.PUT("/user/:id", charupdate.New(ctx, log, db))

	go r.Run(cfg.Server.Addres)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	db.Shutdown()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case local:
		log = slog.New(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case prod:
		log = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case dev:
		log = slog.New(slog.NewJSONHandler(os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}

func MustLoadDSN() string {
	if err := godotenv.Load(); err != nil {
		panic("error loading dsn")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=allow",
		user, password, host, port, dbname,
	)

	return dsn
}
