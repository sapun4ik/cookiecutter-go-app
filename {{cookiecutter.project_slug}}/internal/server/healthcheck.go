package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"{{cookiecutter.module_path}}/pkg/config"
	"{{cookiecutter.module_path}}/pkg/logger"

	"github.com/heptiolabs/healthcheck"
	"github.com/jmoiron/sqlx"
)

type HealthCheckCfg struct {
	Log    logger.Logger
	Cfg    *config.Config
	PsqlDB *sqlx.DB
}

type HealthCheck struct {
	healthCheckCfg *HealthCheckCfg
	httpServer     *http.Server
}

func NewHealthCheck(ctx context.Context, healthCheckCfg *HealthCheckCfg) *HealthCheck {
	health := healthcheck.NewHandler()

	mux := http.NewServeMux()
	mux.HandleFunc(healthCheckCfg.Cfg.Probes.LivenessPath, health.LiveEndpoint)
	mux.HandleFunc(healthCheckCfg.Cfg.Probes.ReadinessPath, health.ReadyEndpoint)

	configureHealthCheckEndpoints(ctx, healthCheckCfg, health)

	return &HealthCheck{
		healthCheckCfg: healthCheckCfg,
		httpServer: &http.Server{
			Handler:           mux,
			Addr:              healthCheckCfg.Cfg.Probes.Port,
			ReadHeaderTimeout: healthCheckCfg.Cfg.HTTP.ReadTimeout,
			WriteTimeout:      healthCheckCfg.Cfg.HTTP.WriteTimeout,
			ReadTimeout:       healthCheckCfg.Cfg.HTTP.ReadTimeout,
		},
	}
}

func (s *HealthCheck) Run() *HealthCheck {
	go func() {
		s.healthCheckCfg.Log.Infof("(%s) Kubernetes probes listening on port: %s", s.healthCheckCfg.Cfg.Application.Name, s.healthCheckCfg.Cfg.Probes.Port)

		if err := s.httpServer.ListenAndServe(); err != nil {
			s.healthCheckCfg.Log.Errorf("(HealthCheck) err: %v", err)
		}
	}()

	return s
}

func (s *HealthCheck) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func configureHealthCheckEndpoints(ctx context.Context, healthCheck *HealthCheckCfg, health healthcheck.Handler) {
	health.AddLivenessCheck("postgres", healthcheck.AsyncWithContext(ctx, func() error {
		if healthCheck.PsqlDB == nil {
			healthCheck.Log.Warnf("(postgres Liveness Check) err: %v", errors.New("database is nil"))

			return errors.New("database is nil")
		}

		return healthCheck.PsqlDB.PingContext(ctx)
	}, time.Duration(healthCheck.Cfg.Probes.CheckIntervalSeconds)*time.Second))

	health.AddReadinessCheck("postgres", healthcheck.AsyncWithContext(ctx, func() error {
		if healthCheck.PsqlDB == nil {
			healthCheck.Log.Warnf("(postgres Liveness Check) err: %v", errors.New("database is nil"))

			return errors.New("database is nil")
		}

		return healthCheck.PsqlDB.PingContext(ctx)
	}, time.Duration(healthCheck.Cfg.Probes.CheckIntervalSeconds)*time.Second))
}
