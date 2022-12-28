package crud

import (
	"github.com/gin-gonic/gin"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/config"

	fn "github.com/mrsimonemms/openfaas-templates/template/golang-crud/function"
)

type handler struct {
	Config *config.Config
}

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	h := &handler{
		Config: cfg,
	}

	routes := r.Group(cfg.RoutePrefix)
	routes.Use(fn.Middleware...)
	if cfg.GetMany {
		routes.GET("/", h.GetMany)
	}
	if cfg.Create {
		routes.POST("/", h.Create)
	}
	if cfg.GetOne {
		routes.GET("/:id", h.GetOne)
	}
	if cfg.Update {
		routes.PATCH("/:id", h.Update)
	}
	if cfg.Delete {
		routes.DELETE("/:id", h.Delete)
	}
}
