package middlewares

import (
	"strings"
	"time"

	"{{cookiecutter.module_path}}/pkg/config"
	"{{cookiecutter.module_path}}/pkg/constants"
	"{{cookiecutter.module_path}}/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MiddlewareMetricsCb func(err error)

type MiddlewareManager interface {
	RequestLoggerMiddleware() gin.HandlerFunc
}

type MiddlewareManagerCfg struct {
	log       logger.Logger
	cfg       *config.Config
	metricsCb MiddlewareMetricsCb
}

func NewMiddlewareManager(log logger.Logger, cfg *config.Config, metricsCb MiddlewareMetricsCb) *MiddlewareManagerCfg {
	return &MiddlewareManagerCfg{log: log, cfg: cfg, metricsCb: metricsCb}
}

func (mw *MiddlewareManagerCfg) RequestLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		// before request

		ctx.Next()

		// after request
		req := ctx.Request

		size := int64(ctx.Writer.Size())
		status := ctx.Writer.Status()
		latency := time.Since(start)

		err := ctx.Errors.Last().Err

		if !mw.checkIgnoredURI(ctx.Request.RequestURI, mw.cfg.HTTP.IgnoreLogUrls) {
			mw.httpMiddlewareAccessLogger(req.Method, req.URL.String(), status, size, latency)
		}

		mw.metricsCb(err)
	}
}

func (mw *MiddlewareManagerCfg) httpMiddlewareAccessLogger(method, uri string, status int, size int64, time time.Duration) {
	mw.log.Info(
		constants.HTTP,
		zap.String(constants.METHOD, method),
		zap.String(constants.URI, uri),
		zap.Int(constants.STATUS, status),
		zap.Int64(constants.SIZE, size),
		zap.Duration(constants.TIME, time),
	)
}

func (mw *MiddlewareManagerCfg) checkIgnoredURI(requestURI string, uriList []string) bool {
	for _, s := range uriList {
		if strings.Contains(requestURI, s) {
			return true
		}
	}

	return false
}
