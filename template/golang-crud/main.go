package main

import (
	"fmt"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	logger "github.com/mrsimonemms/gin-structured-logger"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/config"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/crud"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo/options"

	fn "github.com/mrsimonemms/openfaas-templates/template/golang-crud/function"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := mgm.SetDefaultConfig(nil, cfg.MongoDB.DBName, options.Client().ApplyURI(cfg.MongoDB.URL)); err != nil {
		panic(err)
	}

	// Set log level
	zerolog.SetGlobalLevel(cfg.Logger.Level)

	r := gin.New()
	r.Use(
		requestid.New(requestid.WithCustomHeaderStrKey(cfg.Logger.RequestIDHeader)),
		logger.New(),
		gin.Recovery(),
		func(ctx *gin.Context) {
			logger.Get(ctx).Debug().Str("path", ctx.Request.URL.Path).Msg("New HTTP call")
		},
	)

	crud.RegisterRoutes(r, cfg)

	// Register additional routes
	fn.AdditionalRoutes(r, cfg)

	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
