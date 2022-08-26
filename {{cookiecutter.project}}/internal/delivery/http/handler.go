package http

import (
	"net/http"

	"{{cookiecutter.module_path}}/internal/metrics"
	"{{cookiecutter.module_path}}/pkg/config"
	"{{cookiecutter.module_path}}/pkg/logger"
	"{{cookiecutter.module_path}}/pkg/tracing"

	// import swagger doc
	_ "{{cookiecutter.module_path}}/docs"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	log     logger.Logger
	metrics *metrics.SearchMicroserviceMetrics
}

func NewHandler(log logger.Logger, metrics *metrics.SearchMicroserviceMetrics) *Handler {
	return &Handler{log: log, metrics: metrics}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(
		corsMiddleware,
	)

	if cfg.Application.Profile != config.Prod {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.GET("/ping", h.Ping)

	return router
}

// Ping			 godoc
// @Summary      ping
// @Description  do ping
// @Tags         ping
// @Accept       json
// @Produce      plain
// @Success      200  {string}  string  "pong"
// @Failure      400  {string}  string  "ok"
// @Failure      404  {string}  string  "ok"
// @Failure      500  {string}  string  "ok"
// @Router       /ping [get]
func (h *Handler) Ping(ctx *gin.Context) {
	_, span := tracing.StartHTTPServerTracerSpan(ctx, "pingPong")
	defer span.Finish()
	h.metrics.HTTPSuccessPingRequests.Inc()
	ctx.String(http.StatusOK, "pong")
	span.LogFields(log.String("ping", "pong"))
}
