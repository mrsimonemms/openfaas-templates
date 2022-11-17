package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kamva/mgm/v3"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/config"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/crud"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := mgm.SetDefaultConfig(nil, cfg.MongoDB.DBName, options.Client().ApplyURI(cfg.MongoDB.URL)); err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	crud.RegisterRoutes(r, cfg)

	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
