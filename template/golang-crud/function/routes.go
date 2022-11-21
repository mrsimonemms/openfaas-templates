package function

import (
	"github.com/gin-gonic/gin"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/config"
)

// AdditionalRoutes to be added. These are normal Gin components
// @link https://gin-gonic.com/docs/examples/grouping-routes
var AdditionalRoutes func(r *gin.Engine, cfg *config.Config) = func(r *gin.Engine, cfg *config.Config) {}
