package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	delivery "{{cookiecutter.module_path}}/internal/delivery/http"
	"{{cookiecutter.module_path}}/internal/metrics"
	"{{cookiecutter.module_path}}/internal/repository"
	"{{cookiecutter.module_path}}/internal/server"
	"{{cookiecutter.module_path}}/internal/service"
	"{{cookiecutter.module_path}}/pkg/config"
	"{{cookiecutter.module_path}}/pkg/logger"
	"{{cookiecutter.module_path}}/pkg/middlewares"
	"{{cookiecutter.module_path}}/pkg/postgres"
	"{{cookiecutter.module_path}}/pkg/tracing"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

const (
	waitShotDownDuration = 3 * time.Second
)

type App struct {
	log               logger.Logger
	doneCh            chan struct{}
	cfg               *config.Config
	metricsServer     *server.Metrics
	middlewareManager middlewares.MiddlewareManager
	metrics           metrics.Metrics
	healthCheckServer *server.HealthCheck
	psqlDB            *sqlx.DB
}

func NewApp(log logger.Logger, cfg *config.Config) *App {
	return &App{log: log, doneCh: make(chan struct{}), cfg: cfg}
}

// @title Twino API
// @version 1.0
// @description REST API for Twino

// @host localhost:8000
// @BasePath /api/v1/

// Run initializes whole application.

//nolint:funlen
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	var err error

	a.psqlDB, err = postgres.NewPsqlDB(a.cfg)

	if err != nil {
		a.log.Errorf("Postgresql init: %s", err)

		return err
	}
	defer a.psqlDB.Close()

	a.log.Infof(
		"Postgres connected: host=%s port=%s user=%s dbname=%s connMaxLifetime=%s connMaxIdleTime=%s maxIdleConns=%d maxOpenConns=%d",
		a.cfg.Postgres.Host,
		a.cfg.Postgres.Port,
		a.cfg.Postgres.User,
		a.cfg.Postgres.DBName,
		a.cfg.Postgres.ConnMaxLifetime,
		a.cfg.Postgres.ConnMaxIdleTime,
		a.cfg.Postgres.MaxIdleConns,
		a.cfg.Postgres.MaxOpenConns,
	)

	// enable tracing
	if a.cfg.Jaeger.Enable {
		tracer, closer, err := tracing.NewJaegerTracer(a.cfg.Jaeger)
		if err != nil {
			a.log.Errorf("Error init Jaeger tracing %s", err.Error())

			return err
		}
		defer closer.Close() // nolint: errcheck
		opentracing.SetGlobalTracer(tracer)
	}

	a.middlewareManager = middlewares.NewMiddlewareManager(a.log, a.cfg, a.getHTTPMetricsCb())
	a.metrics = metrics.NewServiceMetrics(a.cfg)

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(a.psqlDB)
	services := service.NewServices(service.Deps{
		Repos: repos,
	})

	handlers := delivery.NewHandler(a.log, a.metrics, services)

	// HTTP Server
	srv := server.NewServer(a.cfg, handlers.Init(a.cfg))

	go func() {
		if err := srv.Run(); err != nil {
			a.log.Errorf("(runHttpServer) err: %v", err)
			cancel()
		}
	}()

	a.log.Infof("%s is listening on PORT: %v", a.cfg.Application.Name, a.cfg.HTTP.Port)

	a.metricsServer = server.NewMetrics(a.cfg, a.log).Run(cancel)
	a.healthCheckServer = server.NewHealthCheck(ctx, &server.HealthCheckCfg{
		Log:    a.log,
		Cfg:    a.cfg,
		PsqlDB: a.psqlDB,
	}).Run()

	<-ctx.Done()

	a.waitShootDown(waitShotDownDuration)

	if err := srv.Stop(ctx); err != nil {
		a.log.Warnf("(Shutdown) err: %v", err)
	}
	if err := a.metricsServer.Stop(ctx); err != nil {
		a.log.Warnf("(Shutdown) metricsServer err: %v", err)
	}
	if err := a.healthCheckServer.Stop(ctx); err != nil {
		a.log.Warnf("(shutDownHealthCheckServer) HealthCheckServer err: %v", err)
	}

	<-a.doneCh

	a.log.Infof("%s app exited properly", a.cfg.Application.Name)

	return nil
}

func (a *App) waitShootDown(duration time.Duration) {
	go func() {
		time.Sleep(duration)
		a.doneCh <- struct{}{}
	}()
}

func (a *App) getHTTPMetricsCb() middlewares.MiddlewareMetricsCb {
	return func(err error) {
		if err != nil {
			a.metrics.ErrorHTTPRequestsInc()
		} else {
			a.metrics.SuccessHTTPRequestsInc()
		}
	}
}
