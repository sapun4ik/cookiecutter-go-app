package v1

import (
	"{{cookiecutter.module_path}}/internal/metrics"
	"{{cookiecutter.module_path}}/internal/service"
	"{{cookiecutter.module_path}}/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	log      logger.Logger
	metrics  metrics.Metrics
	services *service.Services
}

func NewHandler(log logger.Logger, metrics metrics.Metrics, services *service.Services) *Handler {
	return &Handler{
		log:      log,
		services: services,
		metrics:  metrics,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initSampleRoutes(v1)
	}
}
