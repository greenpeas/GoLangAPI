package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	app_interface "seal/internal/app/interface"
	"seal/internal/app/usecase"
	"seal/internal/config"
	"seal/internal/repository/pg"
	"seal/internal/service"
	"seal/internal/service/logger"
	appTransport "seal/internal/transport"
	"seal/internal/transport/rest"
	"seal/internal/transport/rest/middleware"
	pgClient "seal/pkg/client/pg"
	"seal/pkg/validator"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Cfg       *config.Config
	Router    *gin.Engine
	Logger    app_interface.Logger
	Db        *pgxpool.Pool
	JwtWorker appTransport.JwtWorker
	Usecase   *usecase.Usecase
	Validator app_interface.Validator
	Ctx       context.Context
}

func Get(cfg *config.Config) *App {

	ctx := context.Background()

	db := pgClient.NewClient(ctx, cfg.Database.Psql.Dsn)

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

	appValidator := validator.NewValidator()

	return &App{
		cfg,
		gin.New(),
		appLog,
		db,
		appTransport.NewJwtWorker(cfg.Jwt.Secret, cfg.Jwt.Ttl, cfg.Jwt.TtlRefresh),
		usecase.GetUsecase(usecase.Params{
			Ctx:       ctx,
			Db:        db,
			Cfg:       cfg,
			Logger:    appLog,
			Validator: appValidator,
		}),
		appValidator,
		ctx,
	}
}

func (app *App) Clean() {
	defer app.Db.Close()
}

func (app *App) Run() {
	source := "file://" + app.Cfg.Database.Psql.MigrationPath
	pg.RunMigrateUp(source, app.Cfg.Database.Psql.Dsn, app.Logger)

	if app.Cfg.Logger.Level == "debug" {
		app.Router.Use(middleware.Logger(app.Logger))
	}

	if app.Cfg.Router.UseCorsMiddleware {
		app.Router.Use(middleware.CORSMiddleware())
	}

	rest.Register(rest.Handler{
		Usecase:   app.Usecase,
		JwtWorker: app.JwtWorker,
		Logger:    app.Logger,
		Validator: app.Validator,
		Router:    app.Router,
		Cfg:       app.Cfg,
	})

	app.Router.StaticFile("/swagger", "./docs/swagger.json")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Cfg.Router.Port),
		Handler: app.Router,
	}

	stopServices := make(chan bool)

	if app.Cfg.RunCheckTelemetry {
		service.RunCheckTelemetry(service.Params{
			Ctx:      app.Ctx,
			Db:       app.Db,
			Logger:   app.Logger,
			StopChan: stopServices,
		})
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	if app.Cfg.RunCheckTelemetry {
		stopServices <- true
	}

	timeOut, err := time.ParseDuration(app.Cfg.ShutdownTimeout)
	if err != nil {
		app.Logger.Error(err.Error())
		timeOut = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Error("Server Shutdown: ", err)
	}

	<-ctx.Done()
}
