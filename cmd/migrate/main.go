package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"seal/internal/config"
	"seal/internal/repository/pg"
	"seal/internal/service/logger"

	"os"
)

func main() {
	log.Println("Migration module run")

	cfg := config.Get("configs/app.yml")

	appLog, err := logger.GetLogger(logger.Params{
		FilePath:   cfg.Logger.FilePath,
		Level:      cfg.Logger.Level,
		MaxSize:    cfg.Logger.MaxSize,
		MaxBackups: cfg.Logger.MaxBackups,
		MaxAge:     cfg.Logger.MaxAge,
		Compress:   cfg.Logger.Compress,
	})
	if err != nil {
		log.Fatal(err)
	}

	appLog.Info("Migration module run")

	args := os.Args

	source := "file://" + cfg.Database.Psql.MigrationPath
	if len(args) > 2 && args[1] == "create" {
		pg.Create(cfg.Database.Psql.MigrationPath, args[2], appLog)
	} else if len(args) > 1 && args[1] == "down" {
		pg.RunMigrateDown(source, cfg.Database.Psql.Dsn, appLog)
	} else {
		pg.RunMigrateUp(source, cfg.Database.Psql.Dsn, appLog)
	}
}
