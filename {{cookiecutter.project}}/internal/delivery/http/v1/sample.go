package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initSampleRoutes(api *gin.RouterGroup) {
	sample := api.Group("/sample")
	{
		sample.GET("/", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "OK!")
		})
	}
}
