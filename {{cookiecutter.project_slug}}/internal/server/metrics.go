package server

import (
	"context"
	"net/http"

	"{{cookiecutter.module_path}}/pkg/config"
	"{{cookiecutter.module_path}}/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	log        logger.Logger
	httpServer *http.Server
	cfg        *config.Config
}

func NewMetrics(cfg *config.Config, log logger.Logger) *Metrics {
	r := gin.New()
	r.Use(
		gin.Recovery(),
	)
	r.GET(cfg.Probes.PrometheusPath, prometheusHandler())

	return &Metrics{
		log: log,
		cfg: cfg,
		httpServer: &http.Server{
			Addr:              ":" + cfg.Probes.PrometheusPort,
			Handler:           r,
			ReadHeaderTimeout: cfg.HTTP.ReadTimeout,
			ReadTimeout:       cfg.HTTP.ReadTimeout,
			WriteTimeout:      cfg.HTTP.WriteTimeout,
			MaxHeaderBytes:    cfg.HTTP.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Metrics) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Metrics) Run(cancel context.CancelFunc) *Metrics {
	go func() {
		s.log.Infof("Metrics app is running on port: %s", s.cfg.Probes.PrometheusPort)

		if err := s.httpServer.ListenAndServe(); err != nil {
			s.log.Errorf("metricsServer.Start: %v", err)
			cancel()
		}
	}()

	return s
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
